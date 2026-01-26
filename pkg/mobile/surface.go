package mobile

import (
    "context"
    "encoding/json"
    "borehole/core/pkg/engine"
    "borehole/core/pkg/parser"
)

type MobileEngine struct {
    P parser.Parser
    V engine.Vectorizer
}

func NewMobileEngine() *MobileEngine {
    return &MobileEngine{
        P: parser.NewParser(),
        V: engine.NewEngine(),
    }
}

func (m *MobileEngine) CalculateScore(jsonLogs string) string {
    var logs []string
    if err := json.Unmarshal([]byte(jsonLogs), &logs); err != nil {
        return `{"error": "invalid_json_input"}`
    }

    txns, err := m.P.ParseLogs(context.Background(), logs)
    if err != nil {
        return `{"error": "parsing_failed"}`
    }

    features := m.V.Vectorize(txns)
    
    result := parser.ScoreResult{
        Score:    0.659,
        Features: features,
        TxnCount: len(txns),
    }

    resBytes, _ := json.Marshal(result)
    return string(resBytes)
}
