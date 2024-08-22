import { Component, Input } from '@angular/core';
import { FormControl, FormGroup, FormsModule, Validators, ReactiveFormsModule} from '@angular/forms';
import { FloatLabelModule } from 'primeng/floatlabel';
import { DropdownModule } from 'primeng/dropdown';
import { InputTextModule } from 'primeng/inputtext';
import { CardModule } from 'primeng/card';

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
  //Get countries from DB
  loadCountries(): void {
    this.countryService.getCountries().subscribe({
      next: (data) => {
        console.log('Countries:', data);
        // Assurez-vous que `data` est un tableau
        if (Array.isArray(data)) {
          this.countries = data.map((country: any) => ({
            name: country.name.common, 
            code: country.cca2
          }));
          if (this.countries.length > 0) {
            this.selectedCountry = this.countries[0];
            this.loadCities();
          }
        } else {
          console.error('Data format is incorrect', data);
        }
      },
      error: (error) => {
        console.error('Failed to load countries', error);
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
