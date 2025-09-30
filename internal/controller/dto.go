package controller

type createSubscriptionRequestDTO struct {
	UserID      string `json:"user_id" example:"a6fa4d7c-8f90-4f92-912e-92c644c57a1e"`
	ServiceName string `json:"service_name" example:"Yandex Plus"`
	Price       int    `json:"price" example:"1000"`
	StartDate   string `json:"start_date" example:"08-2025"`
	EndDate     string `json:"end_date" example:"09-2025"`
}

type createSubscriptionResponseDTO struct {
	ID string `json:"id" example:"b6fa4d7c-8f90-4f92-912e-92c644c57a1e"`
}

type getTotalPriceResponseDTO struct {
	TotalPrice int `json:"total_price" example:"4600"`
}

type getSubscriptionReadDTO struct {
	ID          string `json:"id" example:"b6fa4d7c-8f90-4f92-912e-92c644c57a1e"`
	UserID      string `json:"user_id" example:"a6fa4d7c-8f90-4f92-912e-92c644c57a1e"`
	ServiceName string `json:"service_name" example:"Yandex Plus"`
	Price       int    `json:"price" example:"1000"`
	StartDate   string `json:"start_date" example:"08-2025"`
	EndDate     string `json:"end_date" example:"09-2025"`
}

type getSubscriptionsResponseDTO struct {
	Subscriptions []getSubscriptionReadDTO `json:"subscriptions"`
}

type updateSubscriptionCreateDTO struct {
	ServiceName string `json:"service_name" example:"Yandex Plus"`
	Price       int    `json:"price" example:"499"`
	StartDate   string `json:"start_date" example:"08-2025"`
	EndDate     string `json:"end_date" example:"09-2025"`
}
