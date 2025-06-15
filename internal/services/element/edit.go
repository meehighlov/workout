package element

import (
	"context"

	"github.com/meehighlov/workout/internal/clients/telegram"
	"github.com/meehighlov/workout/internal/repositories/element"
)

func (s *Service) Edit(ctx context.Context, update *telegram.Update) error {
	keyboard := s.builders.KeyboardBuilder.Keyboard()
	keyboard.AppendAsLine(
		keyboard.NewButton(s.constants.BUTTON_TEXT_NAME, s.builders.CallbackDataBuilder.Build("name", s.constants.COMMAND_EDIT_ELEMENT_REQUEST, "0").String()),
		keyboard.NewButton(s.constants.BUTTON_TEXT_LINK, s.builders.CallbackDataBuilder.Build("link", s.constants.COMMAND_EDIT_ELEMENT_REQUEST, "0").String()),
	)

	params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)
	s.clients.Cache.AppendText(ctx, update.GetChatIdStr(), params.ID)
	s.clients.Cache.SetNextHandler(ctx, update.GetChatIdStr(), s.constants.COMMAND_EDIT_ELEMENT_REQUEST)

	s.clients.Telegram.Reply(ctx, "Что будем редактировать?", update, telegram.WithReplyMurkup(keyboard.Murkup()))

	return nil
}

func (s *Service) EditRequest(ctx context.Context, update *telegram.Update) error {
	s.clients.Telegram.Reply(ctx, "Введите новое значение", update)

	params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)

	nextHandler := ""
	switch params.ID {
	case "name":
		nextHandler = s.constants.COMMAND_EDIT_ELEMENT_NAME_SAVE
	case "link":
		nextHandler = s.constants.COMMAND_EDIT_ELEMENT_LINK_SAVE
	}

	s.clients.Cache.SetNextHandler(ctx, update.GetChatIdStr(), nextHandler)

	return nil
}

func (s *Service) EditNameSave(ctx context.Context, update *telegram.Update) error {
	elementId := s.clients.Cache.GetTexts(ctx, update.GetChatIdStr())[0]

	elementToEdit, err := s.repositories.Element.Get(ctx, &element.Filter{ID: elementId}, nil)
	if err != nil {
		s.clients.Telegram.Reply(ctx, "Возникла непредвиденная ошибка", update)
		return err
	}

	elementToEdit.Name = update.Message.Text

	err = s.repositories.Element.Save(ctx, elementToEdit, nil)
	if err != nil {
		s.clients.Telegram.Reply(ctx, "Возникла непредвиденная ошибка", update)
		return err
	}

	s.clients.Cache.SetNextHandler(ctx, update.GetChatIdStr(), "")

	keyboard := s.builders.KeyboardBuilder.Keyboard()
	backButton := keyboard.NewButton(s.constants.BUTTON_TEXT_BACK, s.builders.CallbackDataBuilder.Build(elementToEdit.ID.String(), s.constants.COMMAND_INFO_ELEMENT, "0").String())
	keyboard.AppendAsLine(backButton)

	s.clients.Telegram.Reply(ctx, "Элемент обновлен", update, telegram.WithReplyMurkup(keyboard.Murkup()))

	return nil
}

func (s *Service) EditLinkSave(ctx context.Context, update *telegram.Update) error {
	elementId := s.clients.Cache.GetTexts(ctx, update.GetChatIdStr())[0]

	elementToEdit, err := s.repositories.Element.Get(ctx, &element.Filter{ID: elementId}, nil)
	if err != nil {
		s.clients.Telegram.Reply(ctx, "Возникла непредвиденная ошибка", update)
		return err
	}

	// todo validate link
	elementToEdit.TutorialLink = update.Message.Text

	err = s.repositories.Element.Save(ctx, elementToEdit, nil)
	if err != nil {
		s.clients.Telegram.Reply(ctx, "Возникла непредвиденная ошибка", update)
		return err
	}

	s.clients.Cache.SetNextHandler(ctx, update.GetChatIdStr(), "")

	keyboard := s.builders.KeyboardBuilder.Keyboard()
	backButton := keyboard.NewButton(s.constants.BUTTON_TEXT_BACK, s.builders.CallbackDataBuilder.Build(elementToEdit.ID.String(), s.constants.COMMAND_INFO_ELEMENT, "0").String())
	keyboard.AppendAsLine(backButton)

	s.clients.Telegram.Reply(ctx, "Элемент обновлен", update, telegram.WithReplyMurkup(keyboard.Murkup()))

	return nil
}
