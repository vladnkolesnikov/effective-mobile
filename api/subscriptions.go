package api

import (
	"effective-mobile/logger"
	"effective-mobile/store"
	"effective-mobile/utils"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type SubscriptionsHandler struct {
	store store.SubscriptionsStore
}

func NewSubscriptionsHandler(store store.SubscriptionsStore) *SubscriptionsHandler {
	return &SubscriptionsHandler{
		store: store,
	}
}

func (sh *SubscriptionsHandler) HandleGetUserSubscriptions(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("userId")
	serviceName := r.URL.Query().Get("serviceName")
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	if err := uuid.Validate(id); err != nil {
		logger.LogError("GetUserSubscriptions: failed to validate user id", err)
		utils.WriteResponse(w, http.StatusBadRequest, utils.Envelope{"error": "invalid user id"})
		return
	}

	if serviceName == "" {
		err := errors.New("serivce name is required")
		logger.LogError("GetUserSubscriptions: serviceName is missing", err)
		utils.WriteResponse(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	rangeStartDate := utils.CustomDate{}
	if err := rangeStartDate.ParseQueryDate(startDate, true); err != nil {
		logger.LogError("GetUserSubscriptions: failed to parse startDate", err)
		utils.WriteResponse(w, http.StatusBadRequest, utils.Envelope{"error": "invalid startDate"})
		return
	}

	rangeEndDate := utils.CustomDate{}
	if err := rangeEndDate.ParseQueryDate(endDate, false); err != nil {
		logger.LogError("GetUserSubscriptions: failed to parse endDate", err)
		utils.WriteResponse(w, http.StatusBadRequest, utils.Envelope{"error": "invalid endDate"})
		return
	}

	userId := uuid.MustParse(id)

	params := &store.RequestParams{
		ServiceName: serviceName,
		StartDate:   rangeStartDate,
		EndDate:     rangeEndDate,
	}

	subscriptions, err := sh.store.GetUserSubscriptions(userId, params)
	if err != nil {
		logger.LogError("GetUserSubscriptions: something went wrong", err)
		utils.WriteResponse(w, http.StatusInternalServerError, utils.Envelope{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	utils.WriteResponse(w, http.StatusOK, utils.Envelope{"subscriptions": subscriptions})
}

func (sh *SubscriptionsHandler) HandleGetTotalSubscriptionCost(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("userId")
	serviceName := r.URL.Query().Get("serviceName")
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	if err := uuid.Validate(id); err != nil {
		logger.LogError("GetTotalSubscriptionCost: failed to validate user id", err)
		utils.WriteResponse(w, http.StatusBadRequest, utils.Envelope{"error": "invalid user id"})
		return
	}

	if serviceName == "" {
		err := errors.New("serivce name is required")
		logger.LogError("GetTotalSubscriptionCost: serviceName is missing", err)
		utils.WriteResponse(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	rangeStartDate := utils.CustomDate{}
	if err := rangeStartDate.ParseQueryDate(startDate, true); err != nil {
		logger.LogError("GetTotalSubscriptionCost: failed to parse startDate", err)
		utils.WriteResponse(w, http.StatusBadRequest, utils.Envelope{"error": "invalid startDate"})
		return
	}

	rangeEndDate := utils.CustomDate{}
	if err := rangeEndDate.ParseQueryDate(endDate, false); err != nil {
		logger.LogError("GetTotalSubscriptionCost: failed to parse endDate", err)
		utils.WriteResponse(w, http.StatusBadRequest, utils.Envelope{"error": "invalid endDate"})
		return
	}

	userId := uuid.MustParse(id)

	params := &store.RequestParams{
		ServiceName: serviceName,
		StartDate:   rangeStartDate,
		EndDate:     rangeEndDate,
	}

	priceAmount, err := sh.store.GetTotalSubscriptionCost(userId, params)
	if err != nil {
		logger.LogError("GetTotalSubscriptionCost: something went wrong", err)
		utils.WriteResponse(w, http.StatusInternalServerError, utils.Envelope{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	utils.WriteResponse(w, http.StatusOK, utils.Envelope{"total_cost": priceAmount})
}

func (sh *SubscriptionsHandler) HandleCreateSubscription(w http.ResponseWriter, r *http.Request) {
	var subscription store.Subscription

	if err := json.NewDecoder(r.Body).Decode(&subscription); err != nil {
		logger.LogError("CreateSubscription: failed to decode body", err)
		utils.WriteResponse(w, http.StatusBadRequest, utils.Envelope{"error": http.StatusText(http.StatusBadRequest)})
		return
	}

	if subscription.StartDate.IsInFuture() {
		err := errors.New("startDate cannot be in the future")
		logger.LogError("CreateSubscription: invalid startDate query value", err)
		utils.WriteResponse(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	if err := sh.store.CreateUserSubscription(&subscription); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23514" {
			logger.LogError("CreateSubscription: failed to create subscription with an expiration date before a start date", err)
			utils.WriteResponse(w, http.StatusBadRequest, utils.Envelope{"error": http.StatusText(http.StatusBadRequest)})
			return
		}

		logger.LogError("CreateSubscription: failed to create subscription", err)
		utils.WriteResponse(w, http.StatusInternalServerError, utils.Envelope{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	utils.WriteResponse(w, http.StatusCreated, utils.Envelope{"subscription": subscription})
}
