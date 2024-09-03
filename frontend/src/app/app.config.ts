import { ApplicationConfig, provideZoneChangeDetection } from '@angular/core';
import { provideRouter } from '@angular/router';
import { MessageService, PrimeNGConfig } from 'primeng/api';

import { routes } from './app.routes';
import { provideClientHydration } from '@angular/platform-browser';
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';
import { provideHttpClient, withFetch } from '@angular/common/http';

export const appConfig: ApplicationConfig = {
  providers: [
    provideZoneChangeDetection({ eventCoalescing: true }), 
    provideRouter(routes), 
    provideClientHydration(), 
    provideAnimationsAsync(),
    provideHttpClient(withFetch()), // Add the missing closing parenthesis
    provideHttpClient(), // Fournit le client HTTP
    MessageService,      // Fournit le service de messages PrimeNG
    PrimeNGConfig,       // Fournit la configuration de PrimeNG
  ]
};