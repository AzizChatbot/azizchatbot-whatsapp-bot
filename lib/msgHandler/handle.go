package msghandler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"

	"github.com/AzizChatbot/azizchatbot-whatsapp-bot/lib/redis"
)

func Handle(client *whatsmeow.Client, messageEvent *events.Message) {
	if messageEvent == nil || messageEvent.Info.IsFromMe || messageEvent.Info.IsGroup {
		return
	}

	AI_URL := os.Getenv("AI_URL")
	userMessage := getMessage(messageEvent)
	userId := messageEvent.Info.Chat.ToNonAD()
	ctx := context.Background()
	redisClient := redis.GetClient()

	if userMessage == "بدأ المحادثة" {
		now := time.Now()
		dailyKey := fmt.Sprintf("session:daily:%s:%s", userId, now.Format("20060102"))
		exists, _ := redisClient.Exists(ctx, dailyKey).Result()
		if exists > 0 {
			sendReply(ctx, client, userId, "لقد قمت ببدء المحادثة اليوم بالفعل، حاول غداً مرة اخرى، او قم بأستخدام موقعنا https://aziz.chat")
			return
		}

		// Set daily expiration to the end of the day
		endOfDay := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		durationUntilEndOfDay := time.Until(endOfDay)
		redisClient.Set(ctx, dailyKey, "true", durationUntilEndOfDay)

		// Set session expiration to 5 minutes
		sessionKey := fmt.Sprintf("session:%s", userId)
		redisClient.HSet(ctx, sessionKey, "count", 0)
		redisClient.Expire(ctx, sessionKey, 5*time.Minute).Err()

		sendReply(ctx, client, userId, "مرحباً بك في المحادثة، يمكنك الآن طرح أسئلتك.")
		return
	}

	sessionKey := fmt.Sprintf("session:%s", userId)
	exists, _ := redisClient.Exists(ctx, sessionKey).Result()
	if exists > 0 {
		messageCount, _ := redisClient.HIncrBy(ctx, sessionKey, "count", 1).Result()
		if messageCount > 3 {
			sendReply(ctx, client, userId, "لقد تجاوزت الحد الأقصى لعدد الرسائل المسموح بها في هذه الجلسة، حاول غداً مرة اخرى، او قم بأستخدام موقعنا https://aziz.chat")
			redisClient.Del(ctx, sessionKey)
			return
		}

		// Process the message
		answer, err := callAI(AI_URL, userMessage)
		if err != nil {
			print("Error calling AI:", err.Error())
			sendReply(ctx, client, userId, "حدث خطأ أثناء معالجة طلبك، حاول مرة اخرى.")
			return
		}
		sendReply(ctx, client, userId, answer)
		return
	}
	// If the session key does not exist, send a message to start the conversation
	sendReply(ctx, client, userId, "مرحباً بك في عزيز المساعد الذكي، لبدأ المحادثة أرسل 'بدأ المحادثة' وسأكون سعيداً بمساعدتك.")
}

func callAI(apiURL, question string) (string, error) {
	requestBody, err := json.Marshal(map[string]string{
		"question": question,
	})
	if err != nil {
		return "", err
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Answer string `json:"answer"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Answer, nil
}

func sendReply(ctx context.Context, client *whatsmeow.Client, userId types.JID, message string) {
	_, err := client.SendMessage(ctx, userId, &waE2E.Message{
		Conversation: proto.String(message),
	})
	if err != nil {
		fmt.Println("Error sending message:", err)
	}
}

func getMessage(messageEvent *events.Message) string {
	if messageEvent == nil {
		return "Error"
	}
	textMessage := messageEvent.Message.GetConversation()
	if textMessage != "" {
		return textMessage
	}
	extendedTextMessage := messageEvent.Message.GetExtendedTextMessage()
	return extendedTextMessage.GetText()

}
