package controller

import (
	"encoding/json"
	_ "github.com/BernsteinMondy/subscription-service/docs"
	"github.com/BernsteinMondy/subscription-service/internal/entity"
	"github.com/google/uuid"
	"github.com/swaggo/http-swagger"
	"net/http"
)

func (c *controller) MapHandlers(mux *http.ServeMux) {
	mux.HandleFunc("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	mux.HandleFunc("GET /subscriptions", c.getSubscriptions)
	mux.HandleFunc("GET /subscriptions/{id}", c.getSubscription)
	mux.HandleFunc("GET /subscriptions/price", c.getSubscriptionsTotalPrice)

	mux.HandleFunc("POST /subscriptions", c.postSubscription)
	mux.HandleFunc("DELETE /subscriptions/{id}", c.deleteSubscription)
	mux.HandleFunc("PUT /subscriptions/{id}", c.putSubscription)
}

// GetSubscriptions godoc
// @Summary Get all subscriptions
// @Description Retrieve all subscriptions
// @Tags subscriptions
// @Produce json
// @Success 200 {object} getSubscriptionsResponseDTO "Array of subscriptions"
// @Failure 500 "Internal Server Error - Returns only status code"
// @Router /subscriptions [get]
func (c *controller) getSubscriptions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	subs, err := c.service.GetAllSubscriptions(ctx)
	if err != nil {
		handleError(w, err)
		return
	}

	subscriptionsResult := make([]getSubscriptionReadDTO, 0, len(subs))
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

	var resp = getSubscriptionsResponseDTO{
		Subscriptions: subscriptionsResult,
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	return
}

// GetSubscriptionsTotalPrice godoc
// @Summary Get total price of subscriptions
// @Description Calculate total price of subscriptions with filtering
// @Tags subscriptions
// @Produce json
// @Param service_name query string false "Filter by Service name"
// @Param user_id query string true "User ID" Format(uuid)
// @Param start_date query string true "Start date (MM-YYYY)"
// @Param end_date query string true "End date (MM-YYYY)"
// @Success 200 {object} getTotalPriceResponseDTO "Total price of all the subscriptions"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /subscriptions/price [get]
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

	var resp = getTotalPriceResponseDTO{
		TotalPrice: int(totalPrice),
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	return
}

// GetSubscription godoc
// @Summary Get subscription by ID
// @Description Retrieve a specific subscription by its ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID" Format(uuid)
// @Success 200 {object} getSubscriptionReadDTO
// @Failure 400 "Bad Request"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /subscriptions/{id} [get]
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

// CreateSubscription godoc
// @Summary Create a new subscription
// @Description Create a new subscription for a user
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body createSubscriptionRequestDTO true "Subscription data"
// @Success 201 {object} createSubscriptionResponseDTO "Returns the ID of the created subscription"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /subscriptions [post]
func (c *controller) postSubscription(w http.ResponseWriter, r *http.Request) {
	var req createSubscriptionRequestDTO

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

	if req.Price < 0 {
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

	var resp = createSubscriptionResponseDTO{
		ID: id.String(),
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	return
}

// DeleteSubscription godoc
// @Summary Delete a subscription
// @Description Cancel and remove a subscription by ID
// @Tags subscriptions
// @Param id path string true "Subscription ID" Format(uuid)
// @Success 200 "OK"
// @Failure 400 "Bad Request"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /subscriptions/{id} [delete]
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

// UpdateSubscription godoc
// @Summary Update a subscription
// @Description Update an existing subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID (UUID)"
// @Param subscription body updateSubscriptionCreateDTO true "Updated subscription data"
// @Success 200 "OK"
// @Failure 400 "Bad Request"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /subscriptions/{id} [put]
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
