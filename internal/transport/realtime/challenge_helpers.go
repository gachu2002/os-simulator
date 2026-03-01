package realtime

import (
	"errors"
	"net/http"

	appchallenges "os-simulator-plan/internal/app/challenges"
)

func respondChallengeServiceError(w http.ResponseWriter, r *http.Request, err error) {
	var svcErr *appchallenges.Error
	if !errors.As(err, &svcErr) {
		respondError(w, r, http.StatusInternalServerError, "internal_error", "internal server error")
		return
	}
	respondError(w, r, svcErr.HTTPStatus, svcErr.Code, svcErr.Message)
}

func toValidatorResultDTO(items []appchallenges.ValidatorResultView) []ValidatorResultDTO {
	out := make([]ValidatorResultDTO, 0, len(items))
	for _, item := range items {
		out = append(out, ValidatorResultDTO{
			Name:     item.Name,
			Type:     item.Type,
			Key:      item.Key,
			Passed:   item.Passed,
			Message:  item.Message,
			Expected: item.Expected,
			Actual:   item.Actual,
		})
	}
	return out
}
