package models

type PredictionRequest struct {
	Data []Comment `json:"data"`
}

type Comment struct {
	Id   int    `json:"id" example:"1"`
	Text string `json:"text" example:"Очень понравилось обслуживание в отделении, но мобильное приложение часто зависает."`
}

type PredictionResponse struct {
	Predictions []Prediction `json:"predictions"`
}

type Prediction struct {
	Id         int      `json:"id" example:"1"`
	Topics     []string `json:"topics" example:"['Обслуживание', 'Мобильное приложение']"`
	Sentiments []string `json:"sentiments" example:"['положительно', 'отрицательно']"`
}
