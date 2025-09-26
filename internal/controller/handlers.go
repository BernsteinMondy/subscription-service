package controller

import (
	"encoding/json"
	"github.com/BernsteinMondy/subscription-service/internal/entity"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (c *controller) MapHandlers(mux *http.ServeMux) {
	mux.HandleFunc("GET /subscriptions", c.getSubscriptions)
	mux.HandleFunc("POST /subscriptions", c.postSubscription)
	mux.HandleFunc("DELETE /subscriptions/{id}", c.deleteSubscription)
}

func (c *controller) postSubscription(w http.ResponseWriter, r *http.Request) {
	var req createSubscriptionDTO

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse(req.StartDate, time.RFC3339)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data := &entity.CreateSubscriptionData{
		UserID:      userID,
		ServiceName: req.ServiceName,
		Price:       int32(req.Price),
		StartDate:   startDate,
	}

	ctx := r.Context()
	id, err := c.service.NewSubscription(ctx, data)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *controller) deleteSubscription(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = c.service.CancelSubscription(ctx, id)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *controller) getSubscriptions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	serviceName := query.Get("service_name")
	userIDStr := query.Get("user_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	filter := &entity.GetSubscriptionsFilter{
		UserID:      userID,
		ServiceName: serviceName,
	}

	ctx := r.Context()
	subscriptions, err := c.service.GetSubscriptionsTotalSumFilter(ctx, filter)
	if err != nil {
		handleError(w, err)
		return
	}

	subscriptionsResult := make([]*getSubscriptionReadDTO, 0, len(subscriptions))
	for _, subscription := range subscriptions {
		subscriptionsResult = append(subscriptionsResult, &getSubscriptionReadDTO{
			ID:          subscription.ID.String(),
			UserID:      subscription.UserID.String(),
			ServiceName: subscription.ServiceName,
			Price:       int(subscription.Price),
			StartDate:   subscription.StartDate.Format(time.RFC3339),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(subscriptionsResult)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
