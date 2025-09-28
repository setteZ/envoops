package client

import (
	"encoding/json"
	"testing"
	"time"

	"env-server/models"
)

// -------------------------------
// Unit test for ParseMessage
// -------------------------------
func TestParseMessage(t *testing.T) {
	tests := []struct {
		name        string
		topic       string
		payload     interface{} // raw struct that will be marshaled into JSON
		expectError bool
		expectNode  *models.NodeData
	}{
		{
			name:  "valid message",
			topic: "nodes/123",
			payload: map[string]interface{}{
				"value":    42.5,
				"quantity": "temperature",
			},
			expectError: false,
			expectNode: &models.NodeData{
				NodeId:   "123",
				Quantity: "temperature",
				Value:    42.5,
			},
		},
		{
			name:        "invalid JSON",
			topic:       "nodes/123",
			payload:     "not-json",
			expectError: true,
		},
		{
			name:  "invalid topic format",
			topic: "nodes", // no node ID
			payload: map[string]interface{}{
				"value":    1,
				"quantity": "humidity",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var payloadBytes []byte
			var err error

			switch v := tt.payload.(type) {
			case string:
				payloadBytes = []byte(v)
			default:
				payloadBytes, err = json.Marshal(v)
				if err != nil {
					t.Fatalf("failed to marshal test payload: %v", err)
				}
			}

			result, err := parseMessage(tt.topic, payloadBytes)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result.NodeId != tt.expectNode.NodeId {
				t.Errorf("NodeId mismatch: got %s, want %s", result.NodeId, tt.expectNode.NodeId)
			}
			if result.Quantity != tt.expectNode.Quantity {
				t.Errorf("Quantity mismatch: got %s, want %s", result.Quantity, tt.expectNode.Quantity)
			}
			if result.Value != tt.expectNode.Value {
				t.Errorf("Value mismatch: got %f, want %f", result.Value, tt.expectNode.Value)
			}

			// Time should parse correctly
			_, err = time.Parse("2006-01-02_15:04:05", result.Time)
			if err != nil {
				t.Errorf("Time not in expected format: %s", result.Time)
			}
		})
	}
}

// -------------------------------
// Mock for database.AddData
// -------------------------------
type mockDatabase struct {
	added []*models.NodeData
}

func (m *mockDatabase) AddData(data *models.NodeData) error {
	m.added = append(m.added, data)
	return nil
}

// -------------------------------
// Integration-style test for onMessageReceived
// -------------------------------
func TestOnMessageReceived(t *testing.T) {
	// fake DB
	mockDB := &mockDatabase{}

	// Save and replace DatabaseAddData
	originalAddData := DatabaseAddData
	DatabaseAddData = mockDB.AddData
	defer func() { DatabaseAddData = originalAddData }()

	tests := []struct {
		name        string
		topic       string
		payload     string
		expectSaved bool
		expectNode  *models.NodeData
	}{
		{
			name:        "valid message",
			topic:       "nodes/456",
			payload:     `{"value": 23.4, "quantity": "temperature"}`,
			expectSaved: true,
			expectNode: &models.NodeData{
				NodeId:   "456",
				Quantity: "temperature",
				Value:    23.4,
			},
		},
		{
			name:        "invalid JSON",
			topic:       "nodes/456",
			payload:     `not-json`,
			expectSaved: false,
		},
		{
			name:        "invalid topic format",
			topic:       "nodes", // missing node id
			payload:     `{"value": 23.4, "quantity": "temperature"}`,
			expectSaved: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// reset DB mock before each subtest
			mockDB.added = nil

			message := fakeMessage{
				topic:   tt.topic,
				payload: []byte(tt.payload),
			}

			onMessageReceived(nil, message)

			if tt.expectSaved {
				if len(mockDB.added) != 1 {
					t.Fatalf("expected 1 data entry, got %d", len(mockDB.added))
				}
				data := mockDB.added[0]

				if data.NodeId != tt.expectNode.NodeId {
					t.Errorf("NodeId mismatch: got %s, want %s", data.NodeId, tt.expectNode.NodeId)
				}
				if data.Quantity != tt.expectNode.Quantity {
					t.Errorf("Quantity mismatch: got %s, want %s", data.Quantity, tt.expectNode.Quantity)
				}
				if data.Value != tt.expectNode.Value {
					t.Errorf("Value mismatch: got %f, want %f", data.Value, tt.expectNode.Value)
				}
			} else {
				if len(mockDB.added) != 0 {
					t.Errorf("expected no data saved, but got %d entries", len(mockDB.added))
				}
			}
		})
	}
}

// -------------------------------
// Fake MQTT.Message for testing
// -------------------------------
type fakeMessage struct {
	topic   string
	payload []byte
}

func (m fakeMessage) Duplicate() bool   { return false }
func (m fakeMessage) Qos() byte         { return 0 }
func (m fakeMessage) Retained() bool    { return false }
func (m fakeMessage) Topic() string     { return m.topic }
func (m fakeMessage) MessageID() uint16 { return 1 }
func (m fakeMessage) Payload() []byte   { return m.payload }
func (m fakeMessage) Ack()              {}
