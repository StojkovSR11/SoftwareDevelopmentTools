# Faza izgradnje
FROM golang:1.19 AS build

# Postavi radni direktorijum
WORKDIR /app

# Kopiraj samo fajlove Go modula
COPY go.mod go.sum ./

# Preuzmi Go module (keširano ako go.mod/go.sum nisu promenjeni)
RUN go mod download

# Kopiraj ostatak izvornog koda
COPY . .

# Izgradi aplikaciju
RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

# Konačna faza
FROM alpine:3.14

# Kopiraj binarni fajl iz faze izgradnje
COPY --from=build /main /main

# Izloži port
EXPOSE 8000

# Pokreni binarni fajl
CMD ["/main"]
