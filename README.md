# Bell Scheduler

A web-based application for managing bell schedules on a Raspberry Pi. The application allows users to create, manage, and trigger bell schedules through a user-friendly interface.

## Features

- User authentication and authorization
- Advanced schedule management with multiple time slots
- Schedule management system with three types:
  - Default schedules: Automatically activated at midnight each day
  - Active schedules: Currently running schedule that controls bell ringing
  - Temporary schedules: Active for the current day only, reset at midnight
- Global configurable bell ring duration
- Real-time bell triggering
- Database migration system
- Responsive web interface

## Hardware Requirements

### Required Components
1. Raspberry Pi (any model with GPIO pins)
2. 5V Relay Module (e.g., [SainSmart 2-Channel Relay Module](https://www.sainsmart.com/products/2-channel-relay-module))
3. Bell System (compatible with relay switching)
4. Power Supply for the bell system
5. Micro USB power supply for the Raspberry Pi

### Wiring Instructions

1. Relay Module Connection:
   ```
   Relay Module    Raspberry Pi
   VCC     ----->  5V (Pin 2)
   GND     ----->  GND (Pin 6)
   IN      ----->  GPIO17 (Pin 11) - Default pin, configurable in settings
   ```

2. Bell System Connection:
   ```
   Bell System    Relay Module
   Live     ----->  COM (Common)
   Neutral  ----->  NO (Normally Open)
   ```

3. Power Supply:
   - Connect the bell system's power supply to the relay module's power input
   - Ensure the power supply voltage matches your bell system requirements
   - The relay module should be powered by the Raspberry Pi's 5V supply

### Safety Considerations

1. Power Supply:
   - Use appropriate power supply for your bell system
   - Ensure all connections are properly insulated
   - Consider using a fuse for additional protection

2. Wiring:
   - Double-check all connections before powering on
   - Keep high-voltage wiring separate from low-voltage wiring
   - Use appropriate wire gauge for your bell system

3. Environment:
   - Keep the system in a dry, protected location
   - Ensure proper ventilation
   - Consider using a protective enclosure

## Software Setup

### Prerequisites

1. Backend:
   - Go 1.17 or later
   - SQLite3
   - Raspberry Pi OS (or compatible Linux distribution)

2. Frontend:
   - Node.js 14.x or later
   - npm or yarn

Install on the Raspberry Pi:
   ```bash
   sudo apt install golang sqlite3 nodejs npm
   ```

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/bell_scheduler.git
   cd bell_scheduler
   ```

2. Backend Setup:
   ```bash
   cd backend
   cp .env.example .env
   # Edit .env with your configuration
   make deps
   make build
   ```

3. Frontend Setup:
   ```bash
   cd frontend
   cp .env.example .env
   # Edit .env with your configuration
   make install
   make build
   ```

4. Database Setup:
   ```bash
   cd backend
   make migrate
   ```

   The project includes a migration system for managing database schema changes. You can use the following commands:
   ```bash
   # Create a new migration
   go run cmd/migrate/main.go create [migration_name]
   
   # Apply all pending migrations
   go run cmd/migrate/main.go up
   
   # Rollback the last migration
   go run cmd/migrate/main.go down
   
   # Show migration status
   go run cmd/migrate/main.go status
   ```

### Configuration

1. Backend Environment Variables:
   ```
   JWT_SECRET=your_jwt_secret
   DB_PATH=./data/bell_scheduler.db
   GPIO_PIN=17
   RING_DURATION=5s
   ```

2. Frontend Environment Variables:
   ```
   VUE_APP_API_URL=http://localhost:8080/api
   VUE_APP_AUTH_TOKEN_KEY=your_token_key
   VUE_APP_AUTH_USER_KEY=your_user_key
   ```

### Running the Application

1. Development Mode:
   ```bash
   # Terminal 1 - Backend
   cd backend
   make run

   # Terminal 2 - Frontend
   cd frontend
   make serve
   ```

2. Production Mode:
   ```bash
   # Build both frontend and backend
   make build

   # Run the application
   make run
   ```

3. Running as a Systemd Service:
   
   To run the Bell Scheduler as a systemd service that starts automatically at boot:
   
   ```bash
   # Follow the instructions in the systemd setup guide
   cat systemd_setup.md
   ```
   
   See [systemd_setup.md](systemd_setup.md) for detailed instructions on setting up the Bell Scheduler as a systemd service.

## Testing

1. Backend Tests:
   ```bash
   cd backend
   make test
   make test-coverage  # For coverage report
   ```

2. Frontend Tests:
   ```bash
   cd frontend
   make test
   make test-coverage  # For coverage report
   ```

## Maintenance

1. Logs:
   - Backend logs are written to stdout/stderr
   - Use systemd journal or log rotation for production

2. Database:
   - Regular backups recommended
   - SQLite database located at `backend/data/bell_scheduler.db`

3. GPIO:
   - Monitor GPIO pin status
   - Check relay module operation
   - Verify bell system functionality

## Troubleshooting

1. GPIO Issues:
   - Check wiring connections
   - Verify GPIO permissions
   - Test relay module with simple script

2. Bell Not Triggering:
   - Verify schedule times
   - Check relay module status
   - Test bell system power

3. Application Issues:
   - Check logs for errors
   - Verify configuration
   - Restart application

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.