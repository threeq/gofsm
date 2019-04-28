package gofsm

import "testing"

func Test_encode(t *testing.T) {
	type args struct {
		src string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"basic 1",
			args{`
				Alice -> Bob: Authentication Request
				Bob --> Alice: Authentication Response
				`},
			"UDhYukJav7JCoKnELT2rKt3AJx9IS2mjoKZDAybCJYp9pCzJ24ejB4qjBW4hTCfFKj3LjLC0Qy2YihWWFwyu5QmK4000__-8iHgS"},

		{"basic 2",
			args{`
				@startuml
				Alice -> Bob: Authentication Request
				Bob --> Alice: Authentication Response
				@enduml
				`},
			"UDhYukJav7GeBaaiAYdDpG7p77CoarCLTEqKdFAJh1GSIqioKlDACfCJIpBpynI2KWjBKujBm0gTyfCKT7Nj5C0QiAWiBiZFAqw5s92Qbm8p7n000F__W8yXUG00"},

		{"basic 3",
			args{`
				@startuml
				Bob -> Alice : hello
				@enduml
				`},
			"UDhYukJav7GeBaaiAYdDpG7pdFAJ57Jj51npCfDJ5QmKCb9pSl8XobBpKc2A00400F__WByDP000"},

		{"state graph",
			args{`
					@startuml
					Title [aa] State Graph
				
					[*] --> graph
					State "aa" as graph {
						A --> B
					}
					graph --> [*]
					@enduml
				`},
			"UDhYuWG1X-AInAAIqjmS23SaioGdLI4wCJ5M8RWaiIHLmRqeiI03B0TH4AqLgw2hQwUG3XVdX2XKIanKKaWiXaWeL4EaE1t1YfqWl5e81L414e_MYeMw8ZKl1UO6G0000F__N8SWnW00"},

		{"中文 uml",
			args{`
					@startuml
					Title [你好] State Graph
				
					[*] --> graph
					State "状态图" as graph {
						A: 我
						B: 你
						A --> B
					}
					graph --> [*]
					@enduml
				`},
			"UDhYuWG1X-AInAAIqjmS23SaioGdLI7woTu5JvVkZLK8BaaiILLmBqeio03BGnH5QyKgwEhQAQJ3nJaX2fMUTsrxrj3uTFO-9ON4OeYyGZL41QUZbSApZebGZfELmfEz2s0oODES8BnQ20NH2nAFreg5EZKrBmNcK400003__tN1C0q0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encode(tt.args.src); got != tt.want {
				t.Errorf("encode() = %v, want %v", got, tt.want)
			}
		})
	}
}
