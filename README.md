# BiasharaTrack: Empowering African MSMEs
## A Comprehensive Business Management Solution

### Overview
BiasharaTrack is an enterprise-grade, offline-first business management platform engineered specifically for African MSMEs. Built with TypeScript and React, it delivers robust inventory management, sales tracking, and financial analytics with a focus on reliability and user experience.

### Quick Start Guide

#### Prerequisites
- Go >= 1.16
- Node.js >= 16.0.0
- ngrok
- Make
- npm
- jq (for parsing ngrok output)

#### Installation Steps
1. Clone the repository:
```bash
git clone https://github.com/Stella-Achar-Oiro/biashara-track.git
cd biashara-track
```

2. Install dependencies:
```bash
make install
```

3. Run the application:
```bash
make run
```

#### Available Make Commands
```bash
make help                 # Show all available commands
make build-backend       # Build the Go backend
make build-frontend      # Build the frontend
make build              # Build both backend and frontend
make run-backend        # Run the Go backend
make run-frontend       # Run the frontend development server
make run               # Run both backend and frontend with ngrok
make clean             # Clean up build artifacts
make install           # Install all dependencies
make start-ngrok       # Start ngrok and update callback URL
```

### Technical Architecture

#### Frontend Architecture
- **Framework**: React 18 with TypeScript
- **State Management**: Context API with custom hooks
- **Styling**: Tailwind CSS with custom design tokens
- **Component Architecture**: Atomic design pattern
- **Internationalization**: React-i18next

### Core Features

#### Dashboard System
```typescript
interface DashboardLayoutProps {
  children: React.ReactNode;
}
```
- Responsive layout with collapsible sidebar
- Real-time notifications system
- Session management
- Mobile-first design approach

#### Inventory Management
```typescript
interface Product {
  id: number;
  name: string;
  barcode?: string;
  quantity: number;
  price: number;
  low_stock_threshold: number;
}
```
- Real-time stock tracking
- Barcode integration
- Low stock alerts
- Image management
- CSV import/export

#### Sales Processing
```typescript
interface Sale {
  id: string;
  products: SaleProduct[];
  paymentMethod: 'cash' | 'mpesa' | 'credit';
  totalAmount: number;
  customerDetails?: CustomerDetails;
}
```
- M-PESA integration
- Credit sales handling
- Receipt generation
- Real-time product search
- Keyboard shortcuts

#### Credit Management
```typescript
interface Credit {
  id: string;
  customerName: string;
  amount: number;
  dueDate: Date;
  status: 'active' | 'overdue' | 'paid';
}
```
- Customer credit tracking
- Payment scheduling
- Automated reminders
- Risk assessment

### Technical Implementation

#### Performance Optimizations
- Debounced search
- Lazy loading components
- Memoized calculations
- Optimistic UI updates
- Service worker implementation

#### Security Measures
- Input sanitization
- XSS prevention
- CSRF protection
- Rate limiting
- Secure session management

#### API Integration
```typescript
interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}
```
- RESTful endpoints
- Error handling
- Rate limiting
- Request validation
- Response caching

### Development Workflow

1. **Starting the Development Environment**
   ```bash
   make run
   ```
   This command will:
   - Start ngrok for M-PESA callback URL
   - Update the `.env` file with the new ngrok URL
   - Run the backend server
   - Start the frontend development server

2. **Backend Development**
   ```bash
   make run-backend
   ```
   The backend will start on `http://localhost:8080`

3. **Frontend Development**
   ```bash
   make run-frontend
   ```
   The frontend will start on `http://localhost:5173`

4. **Cleaning Up**
   ```bash
   make clean
   ```
   This will remove build artifacts and dependencies

### Accessing the Application

- Frontend: `http://localhost:5173`
- Backend API: `http://localhost:8080`
- Ngrok Dashboard: `http://localhost:4040`

### Important Notes

1. **M-PESA Integration**:
   - The application automatically manages ngrok tunneling for M-PESA callbacks
   - Check the ngrok dashboard for tunnel status and requests

2. **Database Connection**:
   - The application uses AWS RDS for MySQL database

### Future Roadmap
- [ ] Advanced analytics
- [ ] Inventory forecasting
- [ ] Multi-branch support
- [ ] Advanced reporting
- [ ] Mobile application

### Contribution Guidelines
1. Fork the repository
2. Create feature branch
3. Commit using conventional commits
4. Submit pull request with comprehensive description

### License
MIT License [LICENSE](LICENSE)

### Contributors
1. [@steoiro](https://learn.zone01kisumu.ke/git/steoiro)

2. [@weakinyi](https://learn.zone01kisumu.ke/git/weakinyi)

3. [@vinomondi](https://learn.zone01kisumu.ke/git/vinomondi)

4. [rcaleb](https://learn.zone01kisumu.ke/git/rcaleb)

5. [@seodhiambo](https://learn.zone01kisumu.ke/git/seodhiambo)
