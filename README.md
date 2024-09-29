# Simple Multiplayer FPS Game WebSocket Server

## Overview
This project is a WebSocket server built with Go for a simple multiplayer first-person shooter (FPS) game. It supports real-time player movement, shooting mechanics, and damage handling, allowing for interactive gameplay.

## Key Features
- Real-time player movement synchronization
- Shooting mechanics with damage calculation
- Efficient WebSocket communication

## Technologies Used
- Go
- Gorilla WebSocket package

## Getting Started

### Prerequisites
- Go (v1.16 or higher)

### Installation

1. Clone the repository.
2. Build the server using Go.
3. Run the server.

The server will run on `http://localhost:8080` by default.

## Usage
Players can connect via WebSocket to send movement and shooting data, which the server processes and broadcasts to all clients.

## Project Structure
```
multiplayer-fps-server/
├── main.go             # Main server file
├── go.mod              # Go module file
└── README.md           # Project documentation
```

## Future Improvements
- Implement advanced game mechanics
- Add user authentication and player stats
- Develop a client-side interface

## License
This project is licensed under the MIT License.

## Acknowledgments
Inspired by various multiplayer game tutorials and Go documentation.

Feel free to contribute or report issues!
