# Groupie Trackers Data

Groupie Trackers Data is a web application that provides information about various artists, their details, and related data.

## Features

- View a list of artists with basic information.
- Access detailed information about each artist, including members, creation date, locations, dates, and relations.
- Responsive design for an optimal viewing experience.

## Getting Started

### Prerequisites

- Go installed on your machine. [Download and Install Go](https://golang.org/dl/)

### Installation

1. Clone the repository:

   git clone https://learn.reboot01.com/git/ahameed/groupie-tracker

2. Navigate to the project directory:

    cd groupie tracker

3. Run the program:

    go run main.go OR go run .

The server will start on http://localhost:2020.

## Usage

- Open your web browser and navigate to http://localhost:2020 to access the Groupie Trackers Data application.
- Explore the list of artists and click on an artist's name to view detailed information.

## Endpoints
- GET /: Displays a list of artists.
- GET /artistData/{artistID}: Displays detailed information about a specific artist.

## Project Structure

groupie-tracker/
│
├── main.go                 # Main application entry point
├── templates/              # HTML templates
│   ├── index.html
│   └── artistData.html
├── static/                 # Static files (CSS, images)
    ├── style.css
    └── concert.jpg
    └── artist.jpg
    └── icon.gif            
├── README.md               # Project documentation
└── go.mod                  # Go module file

# Authors
- Ayman Hameed