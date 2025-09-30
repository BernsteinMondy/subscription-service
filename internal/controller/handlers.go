package controller

import (
	"encoding/json"
	"github.com/BernsteinMondy/subscription-service/internal/entity"
	"github.com/google/uuid"
	"net/http"
)

func (c *controller) MapHandlers(mux *http.ServeMux) {
	mux.HandleFunc("GET /subscriptions", c.getSubscriptions)
	mux.HandleFunc("GET /subscriptions/{id}", c.getSubscription)
	mux.HandleFunc("GET /subscriptions/price", c.getSubscriptionsTotalPrice)

	mux.HandleFunc("POST /subscriptions", c.postSubscription)
	mux.HandleFunc("DELETE /subscriptions/{id}", c.deleteSubscription)
	mux.HandleFunc("PUT /subscriptions/{id}", c.putSubscription)
}

func (c *controller) getSubscriptions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	subs, err := c.service.GetAllSubscriptions(ctx)
	if err != nil {
		handleError(w, err)
		return
	}

	var subscriptionsResult []getSubscriptionReadDTO
	for _, sub := range subs {
		subscriptionsResult = append(subscriptionsResult, getSubscriptionReadDTO{
			ID:          sub.ID.String(),
			UserID:      sub.UserID.String(),
			ServiceName: sub.ServiceName,
			Price:       int(sub.Price),
			StartDate:   sub.StartDate.Format(timeFormat),
			EndDate:     sub.EndDate.Format(timeFormat),
		})
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(map[string][]getSubscriptionReadDTO{
		"subscriptions": subscriptionsResult,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func (c *controller) getSubscriptionsTotalPrice(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	serviceName := query.Get("service_name")
	userIDStr := query.Get("user_id")
	startDateStr := query.Get("start_date")
	endDateStr := query.Get("end_date")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	startDate, endDate, err := parseStartAndEndDate(startDateStr, endDateStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	filter := &entity.GetSubscriptionsFilter{
		UserID:      userID,
		ServiceName: serviceName,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	ctx := r.Context()
	totalPrice, err := c.service.GetSubscriptionsTotalSumFilter(ctx, filter)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"total_price": totalPrice,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func (c *controller) getSubscription(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	sub, err := c.service.GetSubscription(ctx, id)
	if err != nil {
		handleError(w, err)
		return
	}

	var resp = getSubscriptionReadDTO{
		ID:          sub.ID.String(),
		UserID:      sub.UserID.String(),
		ServiceName: sub.ServiceName,
		Price:       int(sub.Price),
		StartDate:   sub.StartDate.Format(timeFormat),
		EndDate:     sub.EndDate.Format(timeFormat),
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func (c *controller) postSubscription(w http.ResponseWriter, r *http.Request) {
	var req createSubscriptionReadDTO

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

	startDate, endDate, err := parseStartAndEndDate(req.StartDate, req.EndDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data := &entity.CreateSubscriptionData{
		UserID:      userID,
		ServiceName: req.ServiceName,
		Price:       int32(req.Price),
		StartDate:   startDate,
		EndDate:     endDate,
	}

	ctx := r.Context()
	id, err := c.service.NewSubscription(ctx, data)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{
		"id": id.String(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	return
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
	return
}

func (c *controller) putSubscription(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req updateSubscriptionCreateDTO

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	startDate, endDate, err := parseStartAndEndDate(req.StartDate, req.EndDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data := &entity.UpdateSubscriptionData{
		ServiceName: req.ServiceName,
		Price:       int32(req.Price),
		StartDate:   startDate,
		EndDate:     endDate,
	}

	ctx := r.Context()
	err = c.service.UpdateSubscription(ctx, id, data)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
