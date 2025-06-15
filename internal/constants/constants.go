package constants

import "github.com/meehighlov/workout/internal/config"

type Constants struct {
	START_MESSAGE string
	ERROR_MESSAGE string

	ELEMENTS_SELECTION_WHEEL_MESSAGE string
	ELEMENTS_LIST_MESSAGE            string
	WORKOUT_SAVED_MESSAGE            string
	WORKOUT_LIST_MESSAGE             string
	WORKOUT_NOT_SAVED_MESSAGE        string
	WORKOUT_DRILL_EXEC_MESSAGE       string

	BUTTON_TEXT_WORKOUT_DRILL_SETS_INCREASE string
	BUTTON_TEXT_WORKOUT_DRILL_SETS_DECREASE string
	BUTTON_TEXT_WORKOUT_DRILL_REPS_INCREASE string
	BUTTON_TEXT_WORKOUT_DRILL_REPS_DECREASE string
	BUTTON_TEXT_ADD                         string
	BUTTON_TEXT_OPEN                        string
	BUTTON_TEXT_BACK                        string
	BUTTON_TEXT_TUTORIAL                    string
	BUTTON_TEXT_EDIT                        string
	BUTTON_TEXT_DELETE                      string
	BUTTON_TEXT_NAME                        string
	BUTTON_TEXT_LINK                        string
	BUTTON_TEXT_SAVE                        string
	BUTTON_TEXT_NEXT                        string
	BUTTON_TEXT_PREV                        string
	BUTTON_TEXT_ERASE_ELEMENT               string
	BUTTON_TEXT_EXEC                        string
	BUTTON_TEXT_DRILLS                      string
	BUTTON_TEXT_WORKOUT_DRILL_WEIGHT_TUNE   string

	COMMAND_START string

	COMMAND_LIST_ELEMENT                   string
	COMMAND_ELEMENT_INFO                   string
	COMMAND_ELEMENTS                       string
	COMMAND_NEW_ELEMENT                    string
	COMMAND_ADD_ELEMENT                    string
	COMMAND_ADD_ELEMENT_SAVE               string
	COMMAND_INFO_ELEMENT                   string
	COMMAND_EDIT_ELEMENT                   string
	COMMAND_EDIT_ELEMENT_SAVE              string
	COMMAND_DELETE_ELEMENT                 string
	COMMAND_ELEMENT_SWITCH_STATUS          string
	COMMAND_EDIT_ELEMENT_NAME_SAVE         string
	COMMAND_EDIT_ELEMENT_LINK_SAVE         string
	COMMAND_EDIT_ELEMENT_REQUEST           string
	COMMAND_DELETE_ELEMENT_CONFIRM         string
	COMMAND_NEW_WORKOUT                    string
	COMMAND_ADD_ELEMENT_TO_WORKOUT         string
	COMMAND_ADD_ELEMENT_TO_WORKOUT_CONTROL string
	COMMAND_SAVE_WORKOUT                   string
	COMMAND_INFO_WORKOUT                   string
	COMMAND_INFO_WORKOUT_PLUS_SETS         string
	COMMAND_INFO_WORKOUT_MINUS_SETS        string
	COMMAND_INFO_WORKOUT_PLUS_REPS         string
	COMMAND_INFO_WORKOUT_MINUS_REPS        string
	COMMAND_WORKOUT_PLUS_SET               string
	COMMAND_WORKOUT_MINUS_SET              string
	COMMAND_WORKOUT_NEXT_SET               string
	COMMAND_WORKOUT_PREV_SET               string
	COMMAND_WORKOUT_PLUS_REPS              string
	COMMAND_WORKOUT_MINUS_REPS             string
	COMMAND_LIST_WORKOUT                   string
	COMMAND_WORKOUTS                       string
	COMMAND_DELETE_WORKOUT                 string
	COMMAND_DELETE_WORKOUT_CONFIRM         string
	COMMAND_EDIT_WORKOUT                   string
	COMMAND_EDIT_WORKOUT_NAME_SAVE         string
	COMMAND_ADD_ELEMENT_TO_WORKOUT_RM_EL   string
	COMMAND_DRILL_EXEC                     string
	COMMAND_EDIT_WORKOUT_DRILLS            string
	COMMAND_WORKOUT_TUNE_WEIGHT            string
}

func New(cfg *config.Config) *Constants {
	return &Constants{
		START_MESSAGE:                           "Привет!",
		ERROR_MESSAGE:                           "Произошла ошибка",
		COMMAND_START:                           "/start",
		COMMAND_LIST_ELEMENT:                    "element_list",
		COMMAND_ELEMENT_INFO:                    "element_info",
		COMMAND_ELEMENTS:                        "/elements",
		COMMAND_NEW_ELEMENT:                     "/new_element",
		BUTTON_TEXT_ADD:                         "➕",
		COMMAND_ADD_ELEMENT:                     "element_add",
		COMMAND_ADD_ELEMENT_SAVE:                "element_add_save",
		BUTTON_TEXT_OPEN:                        "смотреть",
		COMMAND_INFO_ELEMENT:                    "element_info",
		BUTTON_TEXT_BACK:                        "⬅️",
		BUTTON_TEXT_TUTORIAL:                    "туториал",
		BUTTON_TEXT_EDIT:                        "✏️",
		COMMAND_EDIT_ELEMENT:                    "element_edit",
		COMMAND_EDIT_ELEMENT_SAVE:               "element_edit_save",
		BUTTON_TEXT_DELETE:                      "🗑",
		COMMAND_DELETE_ELEMENT:                  "element_delete",
		COMMAND_DELETE_ELEMENT_CONFIRM:          "el_del_cnfrm",
		COMMAND_ELEMENT_SWITCH_STATUS:           "element_switch_status",
		BUTTON_TEXT_NAME:                        "название",
		BUTTON_TEXT_LINK:                        "ссылка",
		COMMAND_EDIT_ELEMENT_NAME_SAVE:          "el_ed_nm_sv",
		COMMAND_EDIT_ELEMENT_LINK_SAVE:          "el_ed_lnk_sv",
		COMMAND_EDIT_ELEMENT_REQUEST:            "el_ed_req",
		COMMAND_ADD_ELEMENT_TO_WORKOUT:          "add_etw",
		COMMAND_ADD_ELEMENT_TO_WORKOUT_CONTROL:  "add_etw_control",
		COMMAND_NEW_WORKOUT:                     "/new_workout",
		ELEMENTS_SELECTION_WHEEL_MESSAGE:        "выберите упражнения\n",
		ELEMENTS_LIST_MESSAGE:                   "Мои элементы",
		BUTTON_TEXT_SAVE:                        "сохранить",
		COMMAND_SAVE_WORKOUT:                    "save_workout",
		WORKOUT_SAVED_MESSAGE:                   "Тренировка сохранена",
		COMMAND_INFO_WORKOUT:                    "workout_info",
		COMMAND_LIST_WORKOUT:                    "workout_list",
		COMMAND_WORKOUTS:                        "/workouts",
		WORKOUT_LIST_MESSAGE:                    "Тренировки",
		BUTTON_TEXT_NEXT:                        "▶️",
		BUTTON_TEXT_PREV:                        "◀️",
		BUTTON_TEXT_WORKOUT_DRILL_SETS_INCREASE: "добавить подход",
		BUTTON_TEXT_WORKOUT_DRILL_SETS_DECREASE: "удалить подход",
		BUTTON_TEXT_WORKOUT_DRILL_REPS_INCREASE: "+ повтор",
		BUTTON_TEXT_WORKOUT_DRILL_REPS_DECREASE: "- повтор",
		COMMAND_INFO_WORKOUT_PLUS_SETS:          "workout_ps",
		COMMAND_INFO_WORKOUT_MINUS_SETS:         "workout_ms",
		COMMAND_INFO_WORKOUT_PLUS_REPS:          "workout_pr",
		COMMAND_INFO_WORKOUT_MINUS_REPS:         "workout_mr",
		COMMAND_DELETE_WORKOUT:                  "workout_delete",
		COMMAND_DELETE_WORKOUT_CONFIRM:          "wo_del_cnfrm",
		COMMAND_EDIT_WORKOUT:                    "workout_edit",
		COMMAND_EDIT_WORKOUT_NAME_SAVE:          "wo_ed_nm_sv",
		COMMAND_WORKOUT_PLUS_SET:                "w_plus_set",
		COMMAND_WORKOUT_MINUS_SET:               "w_minus_set",
		COMMAND_WORKOUT_NEXT_SET:                "w_next_set",
		COMMAND_WORKOUT_PREV_SET:                "w_prev_set",
		COMMAND_WORKOUT_PLUS_REPS:               "w_plus_reps",
		COMMAND_WORKOUT_MINUS_REPS:              "w_minus_reps",
		COMMAND_ADD_ELEMENT_TO_WORKOUT_RM_EL:    "atw_rm_drill",
		BUTTON_TEXT_ERASE_ELEMENT:               "✖️",
		WORKOUT_NOT_SAVED_MESSAGE:               "Тренировка не сохранена - нет упражнений",
		COMMAND_DRILL_EXEC:                      "drill_exec",
		BUTTON_TEXT_EXEC:                        "выполнение",
		BUTTON_TEXT_DRILLS:                      "упражнения",
		COMMAND_EDIT_WORKOUT_DRILLS:             "edit_w_drills",
		BUTTON_TEXT_WORKOUT_DRILL_WEIGHT_TUNE:   "доп вес",
		COMMAND_WORKOUT_TUNE_WEIGHT:             "w_tune_we",
	}
}
