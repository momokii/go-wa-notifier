# Go-WA-Notifier

Welcome to **Go-WA-Notifier**, a fun and innovative web application designed to send the latest news straight to your WhatsApp! With an easy QR code login and seamless API integration with trusted news providers, staying updated has never been more exciting.

## Key Features

- **WhatsApp News Notifications**  
  Get the freshest news delivered directly to your WhatsApp, ensuring you never miss out on the latest buzz.

- **QR Code Login**  
  - Quick and secure authentication via QR code scanning.
  - Connect effortlessly to your WhatsApp account for instant notifications.

- **News API Integration**  
  - Fetch real-time news from trusted news APIs.
  - Enjoy a curated selection of relevant and up-to-date news content.

<!-- - **Notification Dashboard**  
  - **Notification Overview**: View a history of all the news alerts sent to your WhatsApp.
  - **Notification Settings**: Customize the frequency and categories of news you wish to receive.
  - **Account Management**: Easily manage your account preferences and notification options. -->

- **Clean and Modern User Interface**  
  - Built with **HTML**, **Bootstrap**, and **jQuery** for a responsive and engaging user experience.
  - Simple and intuitive navigation that keeps everything within reach.

- **Robust Backend**  
  - Powered by **Golang Fiber** for lightning-fast performance and efficient processing.
  - Scalable architecture designed to handle a high volume of notifications effortlessly.

## Purpose

**Go-WA-Notifier** was created as a fun and practical addition to our portfolio, showcasing the powerful combination of modern web technologies and automated WhatsApp notifications. Whether you're a news junkie or just looking to stay informed, this tool delivers the latest updates directly to you in an engaging and hassle-free way.

## Tech Stack

- **Backend**: Golang Fiber framework
- **Frontend**: HTML, Bootstrap, and jQuery
- **Database**: Postgresql

## Development Status

go-wa-notifier is currently in active development with core functions implemented.


## How to Run

1. **Clone the repository:**
   ```bash
   git clone github.com/momokii/go-wa-notifier
   ```

2. **Install Go:**
   - Make sure Go is installed. You can check by running:
    ```bash
     go version
     ```

3. **Configure Environment Variables**, 
Create a `.env` file in the root directory of the project and fill it with your configuration settings with basic values from `.example.env`.

4. **Run the server:**
     - Start the server using the following command:
     ```bash
     go run main.go
     ```

5. **Optional: Use Air for Hot Reloading:**
     - If you want hot reloading during development, you can use [Air](https://github.com/cosmtrek/air).
   - Start the server with Air by running:
    ```bash
     air
     ```

6. **Access the website:**
     - Open your browser and go to `http://localhost:3004` (or the specified port).
