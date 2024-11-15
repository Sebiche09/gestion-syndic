import { Component, Input, Output, EventEmitter } from '@angular/core';
import { FormGroup, FormsModule, ReactiveFormsModule} from '@angular/forms';
import { FloatLabelModule } from 'primeng/floatlabel';
import { DropdownModule } from 'primeng/dropdown';
import { InputTextModule } from 'primeng/inputtext';
import { CardModule } from 'primeng/card';
import { CommonModule } from '@angular/common';
import { ButtonModule } from 'primeng/button';

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
    CardModule,
    ButtonModule],
  templateUrl: './address.component.html',
  styleUrl: './address.component.scss'
})
export class AddressComponent {
  @Input() addressForm!: FormGroup;
  @Output() next = new EventEmitter<void>();
  @Output() previous = new EventEmitter<void>();

  constructor() {}

  previousStep() {
    this.previous.emit();
  }
  nextStep() {
      this.next.emit();
  }

  isFieldInvalid(field: string): boolean {
    const control = this.addressForm.get(field);
    return !!(control && control.invalid && control.touched);
  }

  getErrorMessage(field: string): string | null {
    const control = this.addressForm.get(field);
    if (control?.errors) {
      if (control.errors['required']) return 'Ce champ est requis.';
    }
    return '';
  }
  
}
