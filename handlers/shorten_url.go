package handlers

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

// Estructura para recibir la petición JSON con la URL a acortar
type ShortenRequest struct {
    URL string `json:"url"`
}

// Handler HTTP para acortar URLs usando la API de CleanURI
func ShortenURL(w http.ResponseWriter, r *http.Request) {
    fmt.Println("ShortenURL Called") // Log para saber que se llamó al handler

    // Decodifica el cuerpo de la petición esperando un JSON con el campo "url"
    var reqData ShortenRequest
    if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    fmt.Printf("URL to shorten: %s\n", reqData.URL) // Log de la URL recibida

    // Prepara el cuerpo del formulario para la petición a CleanURI
    form := []byte("url=" + reqData.URL)
    // Crea la petición HTTP POST a la API de CleanURI
    req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, "https://cleanuri.com/api/v1/shorten", bytes.NewBuffer(form))
    if err != nil {
        http.Error(w, "Failed to create request", http.StatusInternalServerError)
        return
    }
    // Establece el header para indicar que el cuerpo es un formulario
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    // Realiza la petición a CleanURI
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        http.Error(w, "Failed to contact CleanURI", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // Si CleanURI responde con un código diferente a 200 OK, devuelve el error
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        http.Error(w, fmt.Sprintf("CleanURI error: %s", string(body)), http.StatusInternalServerError)
        return
    }

    // Si todo sale bien, copia la respuesta de CleanURI al cliente
    w.Header().Set("Content-Type", "application/json")
    io.Copy(w, resp.Body)
}

// Name: shorten_url
// Inference description: Shortens a long URL using the CleanURI service. Send a JSON body with the "url" field and receive the shortened URL in the response.
// URL: https://related-striking-swift.ngrok-free.app/shorten-url
// Parameters: 
// {
//   "type": "object",
//   "properties": {
//     "url": {
//       "type": "string",
//       "description": "The long URL to be shortened."
//     }
//   },
//   "required": ["url"]
// }
// }
// Return type: String