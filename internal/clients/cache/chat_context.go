package cache

import (
	"context"
	"encoding/json"
	"errors"
)

type WorkoutHolder interface {
	GetID() string
	GetDrills() []string
}

type Workout struct {
	ID     string   `json:"id"`
	Drills []string `json:"drills"`
}

type ChatContext struct {
	ChatId        string   `json:"chatId"`
	UserResponses []string `json:"userResponses"`
	Workout       *Workout `json:"workout"`
	NextHandler   string   `json:"nextHandler"`
}

func newChatContext(chatId string) *ChatContext {
	return &ChatContext{
		ChatId:        chatId,
		UserResponses: []string{},
		NextHandler:   "",
		Workout: &Workout{
			ID:     "",
			Drills: []string{},
		},
	}
}

func (ctx *ChatContext) AppendText(userResponse string) error {
	ctx.UserResponses = append(ctx.UserResponses, userResponse)
	return nil
}

func (ctx *ChatContext) GetTexts() []string {
	return ctx.UserResponses
}

func (ctx *ChatContext) GetNextHandler() string {
	return ctx.NextHandler
}

func (ctx *ChatContext) SetNextHandler(nextHandler string) string {
	ctx.NextHandler = nextHandler
	return ctx.NextHandler
}

func (ctx *ChatContext) reset() error {
	ctx.NextHandler = ""
	ctx.UserResponses = []string{}
	ctx.Workout = &Workout{
		ID:     "",
		Drills: []string{},
	}
	return nil
}

func (c *Client) GetOrCreateChatContext(ctx context.Context, chatId string) *ChatContext {
	val, err := c.Redis.Get(ctx, chatId).Result()

	if err == nil {
		var ctx ChatContext
		if err := json.Unmarshal([]byte(val), &ctx); err == nil {
			return &ctx
		}
	}

	newCtx := newChatContext(chatId)

	jsonCtx, _ := json.Marshal(newCtx)
	c.Redis.Set(ctx, chatId, jsonCtx, c.CacheExpiration)

	return newCtx
}

func (c *Client) saveChatContext(ctx context.Context, chatContext *ChatContext) error {
	jsonCtx, err := json.Marshal(chatContext)
	if err != nil {
		return err
	}

	return c.Redis.Set(ctx, chatContext.ChatId, jsonCtx, c.CacheExpiration).Err()
}

func (c *Client) SetNextHandler(ctx context.Context, chatId string, nextHandler string) error {
	if nextHandler == "" {
		return c.Reset(ctx, chatId)
	}

	chatContext := c.GetOrCreateChatContext(ctx, chatId)
	if chatContext == nil {
		return errors.New("chat context not found")
	}

	chatContext.NextHandler = nextHandler
	return c.saveChatContext(ctx, chatContext)
}

func (c *Client) AppendText(ctx context.Context, chatId string, text string) error {
	chatContext := c.GetOrCreateChatContext(ctx, chatId)
	if chatContext == nil {
		return errors.New("chat context not found")
	}

	chatContext.AppendText(text)
	return c.saveChatContext(ctx, chatContext)
}

func (c *Client) AppendWorkoutElement(ctx context.Context, chatId string, workoutElement string) error {
	chatContext := c.GetOrCreateChatContext(ctx, chatId)
	if chatContext == nil {
		return errors.New("chat context not found")
	}

	if chatContext.Workout == nil {
		chatContext.Workout = &Workout{
			ID:     "",
			Drills: []string{},
		}
	}
	chatContext.Workout.Drills = append(chatContext.Workout.Drills, workoutElement)
	return c.saveChatContext(ctx, chatContext)
}

func (c *Client) PopWorkoutElement(ctx context.Context, chatId string) error {
	chatContext := c.GetOrCreateChatContext(ctx, chatId)
	if chatContext == nil {
		return errors.New("chat context not found")
	}

	if chatContext.Workout == nil {
		return nil
	}

	if len(chatContext.Workout.Drills) == 0 {
		return nil
	}

	chatContext.Workout.Drills = chatContext.Workout.Drills[:len(chatContext.Workout.Drills)-1]
	return c.saveChatContext(ctx, chatContext)
}

func (c *Client) SetWorkoutID(ctx context.Context, chatId string, workoutID string) error {
	chatContext := c.GetOrCreateChatContext(ctx, chatId)
	if chatContext == nil {
		return errors.New("chat context not found")
	}

	if chatContext.Workout == nil {
		chatContext.Workout = &Workout{
			ID:     "",
			Drills: []string{},
		}
	}

	chatContext.Workout.ID = workoutID
	return c.saveChatContext(ctx, chatContext)
}

func (c *Client) GetWorkoutID(ctx context.Context, chatId string) string {
	chatContext := c.GetOrCreateChatContext(ctx, chatId)
	if chatContext == nil {
		return ""
	}

	if chatContext.Workout == nil {
		return ""
	}

	return chatContext.Workout.ID
}

func (c *Client) SetWorkout(ctx context.Context, chatId string, holder WorkoutHolder) error {
	chatContext := c.GetOrCreateChatContext(ctx, chatId)
	if chatContext == nil {
		return errors.New("chat context not found")
	}

	if chatContext.Workout == nil {
		chatContext.Workout = &Workout{
			ID:     "",
			Drills: []string{},
		}
	}

	chatContext.Workout.Drills = holder.GetDrills()
	chatContext.Workout.ID = holder.GetID()
	return c.saveChatContext(ctx, chatContext)
}

func (c *Client) GetWorkoutElements(ctx context.Context, chatId string) []string {
	chatContext := c.GetOrCreateChatContext(ctx, chatId)
	if chatContext == nil {
		return nil
	}

	if chatContext.Workout == nil {
		return nil
	}

	return chatContext.Workout.Drills
}

func (c *Client) GetTexts(ctx context.Context, chatId string) []string {
	chatContext := c.GetOrCreateChatContext(ctx, chatId)
	if chatContext == nil {
		return nil
	}

	return chatContext.GetTexts()
}

func (c *Client) Reset(ctx context.Context, chatId string) error {
	chatContext := c.GetOrCreateChatContext(ctx, chatId)
	if chatContext == nil {
		return errors.New("chat context not found")
	}

	chatContext.reset()
	return c.saveChatContext(ctx, chatContext)
}
