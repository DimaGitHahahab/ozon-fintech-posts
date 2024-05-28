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

type Manager struct {
	schema      graphql.Schema
	subscribers sync.Map
}

func New(schema graphql.Schema) *Manager {
	return &Manager{
		schema: schema,
	}
}

type SubscribeMessage struct {
	Posts []int  `json:"posts"`
	Query string `json:"query"`
}

func (m *Manager) SubscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("failed to upgrade connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	m.handleSubscription(conn)
}

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

type subscriber struct {
	conn          *websocket.Conn
	requestString string
}

func (m *Manager) unsubscribe(subscriptionCancelFn context.CancelFunc, subscriber *subscriber) {
	subscriptionCancelFn()
	if subscriber != nil {
		subscriber.conn.Close()
		m.subscribers.Delete(subscriber)
	}
}

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

		m.manageUnsub(ctx, subscribeChannel, subscriptionCancelFn, sub, msg)
	}()

	return sub
}

func (m *Manager) manageUnsub(
	ctx context.Context, subscribeChannel chan *graphql.Result, subscriptionCancelFn context.CancelFunc,
	sub *subscriber, msg SubscribeMessage,
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
			if err := sendMessage(r, msg, *sub); err != nil {
				if errors.Is(err, websocket.ErrCloseSent) {
					m.unsubscribe(subscriptionCancelFn, sub)
				}
				log.Printf("failed to send message: %v", err)
			}
		}
	}
}

func sendMessage(r *graphql.Result, msg SubscribeMessage, sub subscriber) error {
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
