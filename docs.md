# Eftep - dokumentacja

## Wstęp  

Eftep to prosta i łatwa w użyciu aplikacja do zdalnego zarządzania plikami. Została zaprojektowana z łatwością w użyciu oraz łatwą ścieką do rozwoju protokołu, jednocześnie oferując potężne funkcje do zarządzania dokumentacją.


## Instalacja

### Wymagania

Do zbudowania aplikacji Eftep, wymagane jest posiadanie zainstalowanego środowiska Go w wersji `1.22.3` lub nowszej - zalecana jest instalacja ze strony projektu: https://go.dev/

### Klonowanie repozytorium
```bash
git clone https://github.com/Rakowksiii/eftep.git
```

### Budowanie aplikacji
```bash
cd eftep
go build ./cmd/eftepcli
go build ./cmd/eftepd
```

### Uruchomienie aplikacji
Serwer:
```bash
./eftepd <opcjonalnie nazwa serwera>
```

Client:
```bash
./eftepcli
```

### Przykład działania (Docker)

Przygotowany został plik `docker-compose.yml` pozwalający na uruchomienie serwera oraz klienta w kontenerach Dockerowych. 

Wykonanie:
```bash
docker-compose up
```

spowoduje uruchomienie serwera w 4 kontenerach Dockerowych połączonych siecią. W celu skorzystania z klienta należy podłączyć się do dowolnego kontenera:
```bash
docker exec -it <nazwa kontenera> /bin/bash
```

Następnie uruchomić klienta:
```bash
./eftepcli
```

## Funkcjonalności

### Zarządzanie plikami

Eftep pozwala na zarządzanie plikami w prosty i intuicyjny sposób. Użytkownik może przeglądać, dodawać, usuwać, zmieniać nazwę oraz pobierać pliki z serwera.

#### Korzystanie z klienta

Klient działą w trybie interaktywnym, co oznacza, że użytkownik może korzystać z poleceń w czasie rzeczywistym. W celu uzyskania pomocy, należy wpisać komendę `?`. Każda komenda składa się z pojedynczego ciągu znaków, który jest interpretowany przez klienta. Argumenty są pobierane z kolejnych linii.

##### Konfiguracja
Klient domyślnie korzysta z konfiguracji z pliku `pkg/config/client.go` - definiuje ona lokalizację dla plików, które są pobierane z serwera, grupy multicastowe oraz port na którym klient nasłuchuje podczas wyszykiwania dostępnych serwerów.

#### Korzystanie z serwera

Serwer działa w trybie demonu, co oznacza, że działa w tle i obsługuje żądania klientów. Serwer składa się z dwóch głównych komponentów: `Server` oraz `Discovery`. Pierwszy współbieżnie obsługuje żądania klientów, poprzez połączenia TCP, drugi natomiast nasłuchuje na zdefiniowanych grupach multicastowych w celu informowania klientów o swoim istnieniu.

##### Konfiguracja

Serwer domyślnie korzysta z konfiguracji z pliku `pkg/config/server.go` - definiuje ona lokalizację dla plików, które są przechowywane na serwerze, grupy multicastowe oraz porty na których serwer nasłuchuje. Port nasłuchiwania jest dowolny, ponieważ w procesie `Discovery` serwer wysyła informacje o swojej adresacji. W pliku konfiguracyjnym dodatkowo znajduje się ścieżka do pliku z logami - w katalogu logów systemowych `/var/log/eftepd.log`. 

#### Protokól 

Protokół komunikacji pomiędzy klientem a serwerem oparty jest na protokole binarnym. Protokół jest asynchroniczny, co oznacza, że nagłówek wiadomości kierowanej do serwera różni się od nagłówka wiadomości kierowanej do klienta.

##### Nagłówek klient -> serwer

Nagłówek składa się z 5 bajtów, które są interpretowane w następujący sposób:
- 1 bajt - typ wiadomości
- 4 bajty - długość strumienia danych (big-endian)

Na podstawie typu wiadomości serwer wie, w jaki sposób ma interpretować dane. Dzięki długości strumienia danych serwer wie, ile bajtów ma odczytać, by poprawnie zinterpretować argumenty przesłane przez klienta.

Typy wiadomości:
- `0x00` - pobranie listy plików, 
- `0x01` - pobranie pliku - argumentem jest nazwa pliku
- `0x02` - dodanie pliku - argumentem jest nazwa pliku
- `0x03` - usunięcie pliku - 
- `0x04` - zmiana nazwy pliku

##### Nagłówek serwer -> klient

Nagłówek składa się wyłącznie z 4 bajtów, które są precyzują długość strumienia danych (big-endian). Serwer zawsze odpowiada na żądanie klienta, nawet jeśli żądanie jest niepoprawne. W takim przypadku serwer zwraca informację o błędzie.

### Transmisja plików

Transmisja plików odbywa się w postaci strumieniowej. Po ustaleniu rozmiaru nadciągającego pliku, odbiorca tworzy deskryptor pliku, do którego bezpośrednio zapisuje dane ze strumienia. W ten sposób unikamy konieczności alokacji pamięci na cały plik w pamięci operacyjnej.

### Odkrywanie serwerów

Klient w celu odnalezienia serwerów korzysta z grup multicastowych. Na zdefiniowanych w konfiguracji grupach wysyła wiadomość `ARE_YOU_THERE` w celu odnalezienia serwerów, a następnie oczekuje na odpowiedź przez określony czas. W przypadku braku odpowiedzi, klient zakłada, że nie ma dostępnych serwerów.

Serwer w celu informowania klientów o swoim istnieniu wysyła wiadomość `I_AM_HERE:<nazwa_servera>:<port_tcp>` bezpośrednio do nadawcy. Dzięki temu klient wie o istnieniu serwera, a także o porcie na którym serwer nasłuchuje połączeń.

### Logowanie

Klient z uwagi na swoją interaktywną naturę aktywnie przekazuje informacje o swoim działaniu do użytkownika.

Serwer natomiast zapisuje logi do pliku w katalogu logów systemowych `/var/log/eftepd.log`. Logi zawierają informacje o połączeniach, błędach oraz informacje o operacjach na plikach.


## Rozwój

### Rozszerzalność
Eftep został zaprojektowany z myślą o łatwym rozwoju. W celu dodania nowych funkcjonalności, należy rozszerzyć protokół komunikacji oraz dodać obsługę nowych typów wiadomości w serwerze oraz kliencie.

- 

