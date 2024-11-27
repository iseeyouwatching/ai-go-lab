package main

import (
	"fmt"
	"testing"
)

func TestFactsToStr(t *testing.T) {
	userData := map[string]string{
		"Age":              "25",
		"Favourite colour": "Blue",
	}

	expected := "Age - 25\nFavourite colour - Blue"
	result := factsToStr(userData)

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFactsToStrEmpty(t *testing.T) {
	userData := map[string]string{}

	expected := ""
	result := factsToStr(userData)

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestGetUserState_NewUser(t *testing.T) {
	chatID := int64(12345)

	state := getUserState(chatID)

	if state == nil {
		t.Fatal("Expected non-nil state")
	}
	if state.Stage != "CHOOSING" {
		t.Errorf("Expected stage 'CHOOSING', got '%s'", state.Stage)
	}
	if len(state.UserData) != 0 {
		t.Errorf("Expected empty UserData, got '%v'", state.UserData)
	}
}

func TestGetUserState_ExistingUser(t *testing.T) {
	chatID := int64(12345)

	// Создаем состояние
	state := getUserState(chatID)
	state.Stage = "TYPING_REPLY"
	state.UserData["Age"] = "25"

	// Получаем существующее состояние
	state = getUserState(chatID)
	if state.Stage != "TYPING_REPLY" {
		t.Errorf("Expected stage 'TYPING_REPLY', got '%s'", state.Stage)
	}
	if state.UserData["Age"] != "25" {
		t.Errorf("Expected Age '25', got '%s'", state.UserData["Age"])
	}
}

func TestUserState_Transition(t *testing.T) {
	chatID := int64(12345)
	state := getUserState(chatID)

	// Начальное состояние
	if state.Stage != "CHOOSING" {
		t.Errorf("Expected initial stage 'CHOOSING', got '%s'", state.Stage)
	}

	// Смена стадии
	state.Stage = "TYPING_REPLY"
	if state.Stage != "TYPING_REPLY" {
		t.Errorf("Expected stage 'TYPING_REPLY', got '%s'", state.Stage)
	}

	// Добавление данных пользователя
	state.UserData["Favourite colour"] = "Green"
	if state.UserData["Favourite colour"] != "Green" {
		t.Errorf("Expected Favourite colour 'Green', got '%s'", state.UserData["Favourite colour"])
	}
}

func TestGetUserState(t *testing.T) {
	chatID := int64(12345)

	// Проверяем, что для нового пользователя создаётся состояние
	state := getUserState(chatID)
	if state == nil {
		t.Fatal("Expected non-nil state")
	}
	if state.Stage != "CHOOSING" {
		t.Errorf("Expected stage 'CHOOSING', got '%s'", state.Stage)
	}

	// Проверяем, что возвращается существующее состояние
	state.UserData["Age"] = "25"
	state = getUserState(chatID)
	if state.UserData["Age"] != "25" {
		t.Errorf("Expected '25', got '%s'", state.UserData["Age"])
	}
}

func TestGenerateDoneMessage(t *testing.T) {
	state := &UserState{
		UserData: map[string]string{
			"Age":              "25",
			"Favourite colour": "Blue",
		},
	}

	expected := "I learned these facts about you:\nAge - 25\nFavourite colour - Blue\nUntil next time!"
	result := generateDoneMessage(state)

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	// Тест с пустым состоянием
	state.UserData = map[string]string{}
	expected = "I don't know anything about you yet. Tell me something first!"
	result = generateDoneMessage(state)

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func generateDoneMessage(state *UserState) string {
	if len(state.UserData) == 0 {
		return "I don't know anything about you yet. Tell me something first!"
	}
	return fmt.Sprintf("I learned these facts about you:\n%s\nUntil next time!", factsToStr(state.UserData))
}
