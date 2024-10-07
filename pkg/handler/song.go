package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	musicmax "github.com/MyrzakhmetSmagul/music_max"
)

// getSongs return songs.
//
// @Summary Получить песни из библиотеки
// @Description Получить песни из библиотеки. Можно фильтровать по всем полям кроме id
// @Description структуры. Есть пагинация, если не устанавливать страницу
// @Description и лимит на количество песен будут использованы стандартные значения.
// @Tags Songs
// @Param page query string false "Номер страницы"
// @Param limit query string false "Лимит на количество элементов на одной странице"
// @Param song query string false "Название песни, фильтрация без учета регистра"
// @Param group query string false "Имя исполнителя или название группы, фильтрация без учета регистра"
// @Param link query string false "Ссылка, фильтрация с учетом регистра"
// @Param startDate query string false "Дата выпуска с этой даты включительно"
// @Param endDate query string false "Дата выпуска до этой даты включительно"
// @Produce  json
// @Success 200 {object} musicmax.SongsResponse "Список песен с пагинацией и фильтрами"
// @Failure 400 {object} musicmax.Description "Bad Request"
// @Failure 500 {object} musicmax.Description "Internal Server Error"
// @Router /songs [get]
func (h *Handler) getSongs(w http.ResponseWriter, r *http.Request) {
	slog.Info("GET api/v1/songs")

	page, limit, err := musicmax.GetPageAndLimit(r)

	if err != nil {
		musicmax.DefaultResponse(w, http.StatusBadRequest)
		slog.Warn("error during convertation page or limit value to int",
			slog.String("page", r.URL.Query().Get("page")),
			slog.String("limit", r.URL.Query().Get("limit")),
			slog.Any("error", err))
		return
	}

	slog.Debug("getSongs", slog.Any("page", page), slog.Any("limit", limit))
	slog.Debug("query params", slog.Any("queryParams", r.URL.Query().Encode()))

	filters, err := musicmax.GetQueryParams(r)
	if err != nil {
		musicmax.DefaultResponse(w, http.StatusBadRequest)
		slog.Warn("error during get filters", slog.Any("error", err))
		return
	}

	songs, err := h.service.GetSongs(filters, page, limit)
	if err != nil {
		musicmax.DefaultResponse(w, http.StatusInternalServerError)
		slog.Error("error during convertation \"page\" value to int", slog.String("page", r.URL.Query().Get("page")), slog.Any("error", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(songs)
}

// createSong return song.
//
// @Summary Добавить песню в библиотеку
// @Description Добавить песню в библиотеку
// @Tags Songs
// @Param song body musicmax.SongRequest true "Тело запроса"
// @Produce  json
// @Success 200 {object} musicmax.Description "OK"
// @Failure 400 {object} musicmax.Description "Bad Request"
// @Failure 500 {object} musicmax.Description "Internal Server Error"
// @Router /songs [post]
func (h *Handler) createSong(w http.ResponseWriter, r *http.Request) {
	slog.Info("POST api/v1/songs")

	songReq := new(musicmax.SongRequest)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(songReq)
	if err != nil {
		slog.Warn("handler.createSong error decoding from body", slog.Any("err", err))
		musicmax.DefaultResponse(w, http.StatusBadRequest)
		return
	}

	slog.Debug("song request", slog.Any("songReq", songReq))
	if strings.TrimSpace(songReq.Group) == "" || strings.TrimSpace(songReq.Name) == "" {
		slog.Warn("bad request")
		musicmax.DefaultResponse(w, http.StatusBadRequest)
		return
	}

	err = h.service.CreateSong(songReq)
	if err != nil && !errors.Is(err, musicmax.ErrBadRequest) {
		slog.Error("handler.createSong error creating song", slog.Any("err", err))
		musicmax.DefaultResponse(w, http.StatusInternalServerError)
		return
	} else if errors.Is(err, musicmax.ErrBadRequest) {
		slog.Warn("handler.createSong bad request", slog.Any("err", err))
		musicmax.DefaultResponse(w, http.StatusBadRequest)
		return
	}

	musicmax.DefaultResponse(w, http.StatusOK)
}

// deleteSong delete song.
//
// @Summary Удалить песню из библиотеки
// @Description Удалить песню из библиотеки по id
// @Tags Songs
// @Produce  json
// @Success 200 {object} musicmax.Description "OK"
// @Failure 400 {object} musicmax.Description "Bad Request"
// @Failure 500 {object} musicmax.Description "Internal Server Error"
// @Router /songs/{id} [delete]
func (h *Handler) deleteSong(w http.ResponseWriter, r *http.Request) {
	slog.Info("DELETE /api/v1/songs/{id}")
	id := r.PathValue("id")
	err := h.service.DeleteSong(id)
	if err != nil && !errors.Is(err, musicmax.ErrBadRequest) {
		slog.Error("Handler.deleteSong error", slog.Any("error", err))
		musicmax.DefaultResponse(w, http.StatusInternalServerError)
		return
	}

	if errors.Is(err, musicmax.ErrBadRequest) {
		slog.Warn("Handler.deleteSong bad request", slog.Any("error", err))
		musicmax.DefaultResponse(w, http.StatusBadRequest)
		return
	}

	musicmax.DefaultResponse(w, http.StatusOK)
}

// updateSong update song.
//
// @Summary Измененить данные песни
// @Description Измененить данные песни
// @Tags Songs
// @Param songInfo body musicmax.SongPatchRequest true "Тело запроса"
// @Produce  json
// @Success 200 {object} musicmax.Description "OK"
// @Failure 400 {object} musicmax.Description "Bad Request"
// @Failure 500 {object} musicmax.Description "Internal Server Error"
// @Router /songs/{id} [patch]
func (h *Handler) updateSong(w http.ResponseWriter, r *http.Request) {
	slog.Info("PUT /api/v1/songs/{id}")
	id := r.PathValue("id")
	songReq := new(musicmax.SongPatchRequest)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(songReq)
	if err != nil {
		slog.Warn("handler.createSong error decoding from body", slog.Any("err", err))
		musicmax.DefaultResponse(w, http.StatusBadRequest)
		return
	}
	slog.Debug("song body", slog.Any("body", songReq))

	if songReq.Name == nil &&
		songReq.Group == nil &&
		songReq.Release == nil &&
		songReq.Link == nil &&
		songReq.Text == nil {
		musicmax.DefaultResponse(w, http.StatusBadRequest)
		return
	}

	err = h.service.UpdateSong(id, songReq)
	if err != nil && !errors.Is(err, musicmax.ErrBadRequest) && !errors.Is(err, sql.ErrNoRows) {
		slog.Error("handler.createSong error update song", slog.Any("err", err))
		musicmax.DefaultResponse(w, http.StatusInternalServerError)
		return
	} else if errors.Is(err, musicmax.ErrBadRequest) || errors.Is(err, sql.ErrNoRows) {
		slog.Warn("handler.createSong error update song", slog.Any("err", err))
		musicmax.DefaultResponse(w, http.StatusBadRequest)
		return
	}

	musicmax.DefaultResponse(w, http.StatusOK)
}
