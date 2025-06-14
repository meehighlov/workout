package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"
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
	return &ChatContext{ChatId: chatId, UserResponses: []string{}, NextHandler: ""}
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

func (c *Client) GetOrCreateChatContext(chatId string) *ChatContext {
	ctxChatCreation, cancel := context.WithTimeout(c.ctx, 3*time.Second)
	defer cancel()

	val, err := c.Redis.Get(ctxChatCreation, chatId).Result()

	if err == nil {
		var ctx ChatContext
		if err := json.Unmarshal([]byte(val), &ctx); err == nil {
			return &ctx
		}
	}

	newCtx := newChatContext(chatId)

	jsonCtx, _ := json.Marshal(newCtx)
	c.Redis.Set(ctxChatCreation, chatId, jsonCtx, c.CacheExpiration)

	return newCtx
}

func (c *Client) saveChatContext(ctx *ChatContext) error {
	jsonCtx, err := json.Marshal(ctx)
	if err != nil {
		return err
	}

	ctxSave, cancel := context.WithTimeout(c.ctx, 3*time.Second)
	defer cancel()

	return c.Redis.Set(ctxSave, ctx.ChatId, jsonCtx, c.CacheExpiration).Err()
}

func (c *Client) SetNextHandler(chatId string, nextHandler string) error {
	if nextHandler == "" {
		return c.Reset(chatId)
	}

	ctx := c.GetOrCreateChatContext(chatId)
	if ctx == nil {
		return errors.New("chat context not found")
	}

	ctx.NextHandler = nextHandler
	return c.saveChatContext(ctx)
}

func (c *Client) AppendText(chatId string, text string) error {
	ctx := c.GetOrCreateChatContext(chatId)
	if ctx == nil {
		return errors.New("chat context not found")
	}

	ctx.AppendText(text)
	return c.saveChatContext(ctx)
}

func (c *Client) AppendWorkoutElement(chatId string, workoutElement string) error {
	ctx := c.GetOrCreateChatContext(chatId)
	if ctx == nil {
		return errors.New("chat context not found")
	}

	if ctx.Workout == nil {
		ctx.Workout = &Workout{
			ID:     "",
			Drills: []string{},
		}
	}
	ctx.Workout.Drills = append(ctx.Workout.Drills, workoutElement)
	return c.saveChatContext(ctx)
}

func (c *Client) PopWorkoutElement(chatId string) error {
	ctx := c.GetOrCreateChatContext(chatId)
	if ctx == nil {
		return errors.New("chat context not found")
	}

	if ctx.Workout == nil {
		return nil
	}

	if len(ctx.Workout.Drills) == 0 {
		return nil
	}

	ctx.Workout.Drills = ctx.Workout.Drills[:len(ctx.Workout.Drills)-1]
	return c.saveChatContext(ctx)
}

func (c *Client) SetWorkoutID(chatId string, workoutID string) error {
	ctx := c.GetOrCreateChatContext(chatId)
	if ctx == nil {
		return errors.New("chat context not found")
	}

	ctx.Workout.ID = workoutID
	return c.saveChatContext(ctx)
}

func (c *Client) GetWorkoutID(chatId string) string {
	ctx := c.GetOrCreateChatContext(chatId)
	if ctx == nil {
		return ""
	}

	return ctx.Workout.ID
}

func (c *Client) SetWorkoutElements(chatId string, holder WorkoutHolder) error {
	ctx := c.GetOrCreateChatContext(chatId)
	if ctx == nil {
		return errors.New("chat context not found")
	}

	ctx.Workout.Drills = holder.GetDrills()
	ctx.Workout.ID = holder.GetID()
	return c.saveChatContext(ctx)
}

func (c *Client) GetWorkoutElements(chatId string) []string {
	ctx := c.GetOrCreateChatContext(chatId)
	if ctx == nil {
		return nil
	}

	return ctx.Workout.Drills
}

func (c *Client) GetTexts(chatId string) []string {
	ctx := c.GetOrCreateChatContext(chatId)
	if ctx == nil {
		return nil
	}

	return ctx.GetTexts()
}

func (c *Client) Reset(chatId string) error {
	ctx := c.GetOrCreateChatContext(chatId)
	if ctx == nil {
		return errors.New("chat context not found")
	}

	ctx.reset()
	return c.saveChatContext(ctx)
}
