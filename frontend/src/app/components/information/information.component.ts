import { Component, Input , OnInit} from '@angular/core';
import { FormGroup, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { FloatLabelModule } from 'primeng/floatlabel';
import { InputTextModule } from 'primeng/inputtext';
import { InputTextareaModule } from 'primeng/inputtextarea';
import { CardModule } from 'primeng/card';
import { MessageModule } from 'primeng/message';
import { ButtonModule } from 'primeng/button';
import { DividerModule } from 'primeng/divider';
import { PanelModule } from 'primeng/panel';

@Component({
  selector: 'app-information',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    FormsModule,
    InputTextModule,
    FloatLabelModule,
    InputTextareaModule,
    CardModule,
    MessageModule,
    ButtonModule,
    DividerModule,
    PanelModule
  ],
  templateUrl: './information.component.html',
  styleUrls: ['./information.component.scss']
})
export class InformationComponent {
  @Input() informationForm!: FormGroup;

  // Vérifie si le champ est invalide et a été touché
  isFieldInvalid(field: string): boolean {
    const control = this.informationForm.get(field);
    return !!(control && control.invalid && control.touched);
  }

  // Récupère les messages d'erreur appropriés pour un champ donné
  getErrorMessage(field: string): string | null {
    const control = this.informationForm.get(field);
    if (control?.errors) {
      if (control.errors['required']) return 'Ce champ est requis.';
      if (control.errors['minlength']) return `Minimum de ${control.errors['minlength'].requiredLength} caractères requis.`;
      if (control.errors['maxlength']) return `Maximum ${control.errors['maxlength'].requiredLength} caractères autorisés.`;
      if (control.errors['nameTaken']) return 'Ce nom est déjà pris.';
    }
    return null;
  }
}
