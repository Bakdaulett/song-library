# Song Library API

Welcome to the **Song Library API** project! This API allows you to manage songs, their metadata (group, song title, lyrics), and associated links using a simple RESTful interface.

## Technologies Used
- **Golang**: The backend is built with Go programming language for high performance and concurrency.
- **GIN Framework**: A fast web framework for Go that helps in building APIs.
- **PostgreSQL**: A powerful, open-source relational database used to store song data.
- **Swagger**: Integrated for API documentation, providing interactive API exploration.

## Features

The Song Library API provides the following endpoints:

### 1. **Get All Songs**
- **GET** `/songs`
- **Description**: Fetch all the songs from the database.
- **Response**: Returns a list of songs, including their group, title, release date, and lyrics.

### 2. **Get Song by ID**
- **GET** `/songs/{id}`
- **Description**: Fetch a specific song based on its ID.
- **Response**: Returns details of the song (group, title, release date, lyrics, link).

### 3. **Add New Song**
- **POST** `/songs`
- **Description**: Adds a new song to the database.
- **Request Body Example**:
    ```json
    {
        "group": "Imagine Dragons",
        "song": "Demons",
        "release_date": "2017-02-01T00:00:00Z",
        "lyrics": "When the days are cold and the cards all fold\n\nAnd the saints we see are all made of gold \n\nWhen your dreams all fail and the ones we hail",
        "link": "https://example.com/demons"
    }
    ```
- **Response**: Confirms that the song has been added successfully.

### 4. **Update Song**
- **PUT** `/songs/{id}`
- **Description**: Updates the details of an existing song based on its ID.
- **Request Body Example**:
    ```json
    {
        "group": "Imagine Dragons",
        "song": "Thunder",
        "release_date": "2017-02-01T00:00:00Z",
        "lyrics": "Thunder, feel the thunder \n\nLightning then the thunder \n\nThunder, feel the thunder",
        "link": "https://example.com/thunder"
    }
    ```
- **Response**: Confirms that the song has been updated successfully.

### 5. **Delete Song**
- **DELETE** `/songs/{id}`
- **Description**: Deletes a song from the database by its ID.
- **Response**: Confirms that the song has been deleted successfully.

### **Lyrics API Endpoints**

The following endpoints allow you to manage and retrieve lyrics for songs.

### 6. **Get All Lyrics for a Song**
- **GET** `/songs/{id}/lyrics`
- **Description**: Retrieves all the lyrics for a specific song by its ID.
- **Parameters**:
    - `id`: The unique identifier for the song (integer).
- **Response Example**:
    ```json
    {
        "lyrics": "Verse 1: Have you got color in your cheeks?\n\nVerse 2: Do you ever get that fear that you can’t shift\n\nVerse 3: Crawling back to you..."
    }
    ```

### 7. **Get Specific Lyrics by Verse Number**
- **GET** `/songs/{id}/lyrics/{verse_number}`
- **Description**: Retrieves the lyrics for a specific verse of a song.
- **Parameters**:
    - `id`: The unique identifier for the song (integer).
    - `verse_number`: The verse number you want to retrieve.
- **Response Example** (for `/songs/3/lyrics/1`):
    ```json
    {
        "lyrics": "Verse 1: Have you got color in your cheeks?"
    }
    ```

### 8. **Get Multiple Lyrics by Verse Range**
- **GET** `/songs/{id}/lyrics/{verse_start}-{verse_end}`
- **Description**: Retrieves lyrics for a specific range of verses from a song.
- **Parameters**:
    - `id`: The unique identifier for the song (integer).
    - `verse_start`: The starting verse number.
    - `verse_end`: The ending verse number.
- **Response Example** (for `/songs/3/lyrics/2-3`):
    ```json
    {
        "lyrics": "Verse 2: Do you ever get that fear that you can’t shift\n\nVerse 3: Crawling back to you..."
    }
    ```
  

## API Documentation

Swagger has been integrated into the project for easy API exploration.

- **Swagger URL**: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

You can view the interactive API documentation and test the endpoints directly from your browser.

## Installation

To get started with this project, follow these steps:

### 1. Clone the repository:
```bash
git clone https://github.com/Bakdaulett/song-library.git
cd song-library 
```
### 2. Install dependencies:
Make sure you have Go installed (v1.18 or higher).

Run the following command to install the necessary dependencies:

```bash
go mod tidy
```
### 3. Setup PostgreSQL:
You need to have PostgreSQL installed and running.

Create a new PostgreSQL database (e.g., song_library).
Update the database configuration in the config file (e.g., config/config.go) with your PostgreSQL connection details:
```bash
const (
    DB_HOST     = "localhost"
    DB_PORT     = "5432"
    DB_USER     = "your-username"
    DB_PASSWORD = "your-password"
    DB_NAME     = "song_library"
)
```
### 4. Run the application:
To start the application, simply run:
```bash
go run main.go
```
The API server will be running at http://localhost:8080.

### 5. Run Database Migrations (if needed):
Make sure to run any necessary migrations to create the required tables in your PostgreSQL database.

### 6. View the Swagger UI:
Open your browser and navigate to http://localhost:8080/swagger/index.html to view and interact with the API documentation.

