import { Component, Input } from '@angular/core';
import { FormControl, FormGroup, FormsModule, Validators, ReactiveFormsModule} from '@angular/forms';
import { FloatLabelModule } from 'primeng/floatlabel';
import { DropdownModule } from 'primeng/dropdown';
import { InputTextModule } from 'primeng/inputtext';
import { CardModule } from 'primeng/card';
import { CommonModule } from '@angular/common';
import { CountryService } from '../../services/country.service';
import { CityService } from '../../services/city.services';

interface Country {
  name: string;
  code: string;
}
interface City {
  name: string;
}

@Component({
  selector: 'app-address',
  standalone: true,
  imports: [
    ReactiveFormsModule,
    FormsModule,
    CommonModule,
    InputTextModule,
    FloatLabelModule,
    DropdownModule,
    CardModule],
  templateUrl: './address.component.html',
  styleUrl: './address.component.scss'
})
export class AddressComponent {
  @Input() addressForm!: FormGroup;

  countries: Country[] = [];
  selectedCountry?: Country;

  cities: City[] = [];
  selectedCity?: City; 

  constructor(private countryService: CountryService, private cityService: CityService) {}

  //fonction init
  ngOnInit(): void{
    this.loadCountries();
  }
 // Get countries from DB
loadCountries(): void {
  this.countryService.getCountries().subscribe({
    next: (data) => {
      console.log('Countries:', data);

      // Vérifiez que `data` est bien un tableau et contient des pays valides
      if (Array.isArray(data) && data.length > 0) {
        // Mappez les pays pour récupérer uniquement les informations pertinentes
        this.countries = data.map((country: any) => ({
          name: country.name.common,  // Assurez-vous que `name.common` est correct
          code: country.cca2          // Utilisez `cca2` pour les codes ISO 3166-1 alpha-2
        }));

        // Si vous souhaitez définir un pays par défaut basé sur une logique spécifique
        const defaultCountryCode = 'BE';  // Exemple : 'FR' pour la France
        const defaultCountry = this.countries.find(
          (country) => country.code === defaultCountryCode
        );

        // Si un pays par défaut est trouvé, le sélectionner, sinon ne rien sélectionner
        if (defaultCountry) {
          this.selectedCountry = defaultCountry;
          console.log(`Default country selected: ${this.selectedCountry.name}`);
        } else {
          console.warn('No default country found, please select manually.');
        }

        // Vous pouvez également charger les villes pour ce pays par défaut, si nécessaire
        if (this.selectedCountry) {
          this.loadCities();  // Charge les villes du pays par défaut
        }
      } else {
        console.error('Data format is incorrect or empty:', data);
      }
    },
    error: (error) => {
      console.error('Failed to load countries:', error);
    }
  });
}

  //Get cities from DB
  loadCities(): void {
    if (this.selectedCountry) {
      this.cityService.getCities(this.selectedCountry.code).subscribe({
        next: (data) => {
          console.log('Cities:', data);
          if (data) {      
          } else {
            console.error('Data format is incorrect', data);
          }
        },
        error: (error) => {
          console.error('Failed to load cities', error);
        }
      });
    }
  }
  onCountryChange(event: any): void {
    const selectedCountryCode = event.value; // assuming value contains country code
    this.selectedCountry = this.countries.find(country => country.code === selectedCountryCode);
    this.loadCities();
  }
}
