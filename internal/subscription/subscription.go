package subscription

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/graphql-go/graphql"
)

// Manager handles websocket subscriptions
// Based on: https://github.com/graphql-go/subscription-example
type Manager struct {
	schema      graphql.Schema
	subscribers sync.Map
}

func New(schema graphql.Schema) *Manager {
	return &Manager{
		schema: schema,
	}
}

// SubscribeMessage is the message sent by the client to subscribe to posts
type SubscribeMessage struct {
	Posts []int  `json:"posts"`
	Query string `json:"query"`
}

// SubscriptionsHandler upgrades the connection to a websocket connection
func (m *Manager) SubscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("failed to upgrade connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	m.handleSubscription(conn)
}

// handleSubscription reads the message from the websocket connection
func (m *Manager) handleSubscription(conn *websocket.Conn) {
	subscriptionCtx, subscriptionCancelFn := context.WithCancel(context.Background())
	defer subscriptionCancelFn()

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("failed to read websocket message: %v", err)
			return
		}

		var msg SubscribeMessage
		if err := json.Unmarshal(p, &msg); err != nil {
			log.Printf("failed to unmarshal websocket message: %v", err)
			continue
		}
		m.subscribe(subscriptionCtx, subscriptionCancelFn, conn, msg)
	}
}

// subscriber represents a websocket connection with a query
type subscriber struct {
	conn          *websocket.Conn
	requestString string
}

// unsubscribe closes the connection and removes the subscriber from the list
func (m *Manager) unsubscribe(subscriptionCancelFn context.CancelFunc, subscriber *subscriber) {
	subscriptionCancelFn()
	if subscriber != nil {
		subscriber.conn.Close()
		m.subscribers.Delete(subscriber)
	}
}

// subscribe creates a new subscriber and starts a goroutine to manage the subscription
func (m *Manager) subscribe(ctx context.Context, subscriptionCancelFn context.CancelFunc, conn *websocket.Conn, msg SubscribeMessage) *subscriber {
	sub := &subscriber{
		conn:          conn,
		requestString: msg.Query,
	}
	m.subscribers.Store(&sub, struct{}{})

	ctx = context.WithValue(ctx, "posts", msg.Posts)
	go func() {
		subscribeParams := graphql.Params{
			Context:       ctx,
			RequestString: msg.Query,
			Schema:        m.schema,
		}

		subscribeChannel := graphql.Subscribe(subscribeParams)

		m.manageUnsub(ctx, subscribeChannel, subscriptionCancelFn, sub)
	}()

	return sub
}

// manageUnsub sends the subscription result to the client or unsubscribes if client is disconnected
func (m *Manager) manageUnsub(
	ctx context.Context, subscribeChannel chan *graphql.Result, subscriptionCancelFn context.CancelFunc,
	sub *subscriber,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case r, isOpen := <-subscribeChannel:
			if !isOpen {
				m.unsubscribe(subscriptionCancelFn, sub)
				return
			}
			if err := sendMessage(r, *sub); err != nil {
				if errors.Is(err, websocket.ErrCloseSent) {
					m.unsubscribe(subscriptionCancelFn, sub)
				}
				log.Printf("failed to send message: %v", err)
			}
		}
	}
}

// sendMessage sends the result of the subscription to the client
func sendMessage(r *graphql.Result, sub subscriber) error {
	message, err := json.Marshal(map[string]any{
		"payload": r.Data,
	})
	if err != nil {
		return err
	}

	if err := sub.conn.WriteMessage(websocket.TextMessage, message); err != nil {
		return err
	}

	return nil
}
