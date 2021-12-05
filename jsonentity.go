package main

type ScriptFilterOutput struct {
	Rerun     float32 `json:"rerun"` // 再実行への待機時間. 再実行する場合に設定.
	Variables struct {
		Runned     int    `json:"runned"`       // すでに実行した回数
		Browser    string `json:"browser"`      // ブラウザ名
		BrowserURL string `json:"browserUrl"`   // ブラウザから取得した URL
		Title      string `json:"browserTitle"` // ブラウザから取得した URL
	} `json:"variables"`
	Items []*ScriptFilterItem `json:"items"`
}

type ScriptFilterItem struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle,omitempty"`
	Arg      string `json:"arg"`
	Valid    bool   `json:"valid,omitempty"`
}

func NewScriptFilterItem(title string, subtitle string, arg string, valid bool) *ScriptFilterItem {
	output := ScriptFilterItem{
		Title:    title,
		Subtitle: subtitle,
		Arg:      arg,
		Valid:    valid,
	}
	return &output
}

func (output *ScriptFilterOutput) addItem(item *ScriptFilterItem) {
	output.Items = append(output.Items, item)
}
