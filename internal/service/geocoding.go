package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Este struct solo se usa aquí para mapear la respuesta de la API externa
type GeocodeResponse struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

// GetCoordinates es un método del service para mantener la coherencia
func (s *OrderService) GetCoordinates(address string) (float64, float64, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	
	// Filtramos por Rafaela para que la búsqueda sea precisa
	query := fmt.Sprintf("%s, Rafaela, Argentina", address)
	apiURL := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json&limit=1", url.QueryEscape(query))

	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("User-Agent", "TrackingApp-Zoe-StudentProject") // Requerido por Nominatim

	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var results []GeocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return 0, 0, err
	}

	if len(results) == 0 {
		return 0, 0, errors.New("no se encontró la dirección en Rafaela")
	}

	lat, errLat := strconv.ParseFloat(results[0].Lat, 64)
	lon, errLon := strconv.ParseFloat(results[0].Lon, 64)

	if errLat != nil || errLon != nil {
        return 0, 0, errors.New("formato de coordenadas inválido de la API externa")
    }
	return lat, lon, nil
}