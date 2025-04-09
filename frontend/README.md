# Bell Scheduler Frontend

This is the frontend application for the Bell Scheduler, built with Vue.js 2 and Vuetify.

## Project Structure

```
frontend/
├── src/
│   ├── assets/     # Static assets
│   ├── components/ # Reusable components
│   ├── layouts/    # Layout components
│   ├── plugins/    # Vue plugins
│   ├── router/     # Vue Router configuration
│   ├── store/      # Vuex store modules
│   ├── utils/      # Utility functions
│   ├── views/      # Page components
│   ├── App.vue     # Root component
│   └── main.js     # Application entry point
├── public/         # Public static files
└── package.json    # Project dependencies
```

## Prerequisites

- Node.js 14.x or later
- npm or yarn

## Setup

1. Install dependencies:
   ```bash
   npm install
   ```

2. Create a `.env` file in the root directory:
   ```
   VUE_APP_API_URL=http://localhost:8080/api
   ```

3. Start the development server:
   ```bash
   npm run serve
   ```

4. Build for production:
   ```bash
   npm run build
   ```

## Features

- User authentication and authorization
- Advanced schedule management with support for:
  - Default schedules (automatically activated daily)
  - Active schedules (currently running)
  - Temporary schedules (reset at midnight)
- Time slot configuration with day-specific settings
- Global settings management
- Admin user management
- Dashboard with current schedule status
- Responsive design with Vuetify
- Form validation with vee-validate
- Date/time utilities with moment.js

## Development

The application follows these patterns:
- Component-based architecture
- Vuex for state management
- Vue Router for navigation
- Vuetify for UI components
- Modular store structure
- Reusable components

## Available Scripts

- `npm run serve` - Start development server
- `npm run build` - Build for production
- `npm run lint` - Lint code
- `npm run test:unit` - Run unit tests
- `npm run test:e2e` - Run end-to-end tests