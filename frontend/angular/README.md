# Angular Frontend

Angular 18 standalone component application with enterprise-grade architecture.

## Architecture

This project follows the **core/features/shared** pattern for scalable, maintainable Angular applications:

```
src/app/
├── core/
│   └── services/
│       └── api.service.ts        # Singleton API service (providedIn: 'root')
├── shared/
│   └── components/
│       ├── navbar/
│       │   ├── navbar.component.ts
│       │   ├── navbar.component.html
│       │   └── navbar.component.scss
│       └── footer/
│           ├── footer.component.ts
│           ├── footer.component.html
│           └── footer.component.scss
└── features/
    ├── home/
    │   ├── home.component.ts
    │   ├── home.component.html
    │   ├── home.component.scss
    │   └── home.routes.ts
    └── about/
        ├── about.component.ts
        ├── about.component.html
        ├── about.component.scss
        └── about.routes.ts
```

### Key Concepts

- **Standalone Components**: All components are standalone (no NgModules)
- **Lazy Loading**: Feature modules are loaded on-demand via the router
- **Core Services**: Singleton services provided at root level (e.g., ApiService)
- **Shared Components**: Reusable UI components used across features
- **Feature Modules**: Self-contained feature areas with their own routes

## Development

### Prerequisites

- Node.js 18+
- npm 9+

### Install Dependencies

```bash
npm install
```

### Development Server

```bash
npm start
```

Navigate to `http://localhost:4200/`. The app will automatically reload if you change any of the source files.

### Build

```bash
npm run build
```

Build artifacts will be stored in the `dist/` directory.

## Routing

Routes are defined in `src/app/app.routes.ts` using lazy loading:

```typescript
export const routes: Routes = [
  {
    path: '',
    loadChildren: () => import('./features/home/home.routes').then(m => m.homeRoutes)
  },
  {
    path: 'about',
    loadChildren: () => import('./features/about/about.routes').then(m => m.aboutRoutes)
  }
];
```

Each feature defines its own routes in a `*.routes.ts` file.

## API Integration

The `ApiService` (in `core/services/`) provides centralized API access:

```typescript
@Injectable({ providedIn: 'root' })
export class ApiService {
  private baseUrl = 'http://localhost:8080';

  constructor(private http: HttpClient) {}

  getHealth(): Observable<any> {
    return this.http.get(`${this.baseUrl}/health`);
  }

  getInfo(): Observable<any> {
    return this.http.get(`${this.baseUrl}/info`);
  }
}
```

## Styling

- Global styles: `src/styles.scss`
- Component styles: Each component has an accompanying `.scss` file
- Uses SCSS for preprocessing

## Testing

```bash
npm test
```

## Production

The production build is optimized with:
- Code splitting via lazy loading
- Tree-shaking
- Minification
- Build optimizer

Build with:

```bash
npm run build -- --configuration production
```
