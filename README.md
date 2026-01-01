# Barbershop Telegram Bot

A Telegram bot for barbershops that allows barbers to invite customers, set visit dates, and send periodic reminders.

## Features

- User roles: Barber and Customer
- Invite customers
- Set appointment dates
- View upcoming visits
- Automatic reminders for upcoming visits and to book new appointments

## Setup

1. Create a Telegram bot and get the token from [@BotFather](https://t.me/botfather).

2. Create a `.env` file in the project root with your bot token:
   ```
   TELEGRAM_BOT_TOKEN=your_token_here
   ```

3. Start the bot and send `/start` or `start` to it. The first user to do this will automatically become a barber.

## Running Locally

1. Install Go 1.21 or later.

2. Clone the repository and navigate to the directory.

3. Install dependencies:
   ```bash
   go mod download
   ```

4. Run the bot:
   ```bash
   TELEGRAM_BOT_TOKEN=your_token_here go run .
   ```

5. Run with reflex watcher:
```bash
reflex -r '\.go$' -s -- go run . 
```

## Running with Docker

1. Set the environment variable:
   ```bash
   export TELEGRAM_BOT_TOKEN=your_token_here
   ```

2. Build and run with Docker Compose:
   ```bash
   docker-compose up --build -d
   ```

PostgreSQL will be automatically set up with the database.

## Bot Commands

Commands can be used with or without the `/` prefix.

### For Barbers
- `start` or `/start` - Welcome message
- `invite @username` or `/invite @username` - Invite a customer
- `setvisit @username YYYY-MM-DD HH:MM` or `/setvisit @username YYYY-MM-DD HH:MM` - Set visit date and time
- `makebarber @username` or `/makebarber @username` - Promote user to barber
- `myclients` or `/myclients` - View your clients
- `myvisits` or `/myvisits` - View your visits

### For Admins
- `admins` or `/admins` - View all admins

### For Customers
- `start` or `/start` - Welcome message
- `myvisits` or `/myvisits` - View upcoming appointments
- `createbarber` or `/createbarber` - Become a barber

## Reminders

- Reminds customers 24 hours before their visit
- Reminds customers to book a new appointment if they haven't visited in the last month
- Runs every hour

## Database

Uses PostgreSQL for data storage. Tables:
- `users`: Stores user information and roles
- `visits`: Stores visit dates and statuses
