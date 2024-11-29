import { Component, ViewChild } from '@angular/core';
import { Router } from '@angular/router';
import { DialogModule } from 'primeng/dialog';
import { TableModule } from 'primeng/table';
import { ButtonModule } from 'primeng/button';
import { CondominiumComponent } from '../condominium/condominium.component';
import { ConfirmationService } from 'primeng/api';
import { ConfirmDialogModule } from 'primeng/confirmdialog';
@Component({
    selector: 'displayCondominium',
    templateUrl: 'display-condominium.component.html',
    standalone: true,
    imports: [TableModule, ButtonModule, DialogModule, CondominiumComponent, ConfirmDialogModule],
    providers: [ConfirmationService]
})
export class DisplayCondominiumComponent {
    displayDialog: boolean = false;

    @ViewChild(CondominiumComponent) dialogForm!: CondominiumComponent; 

    constructor(
        private router: Router,
        private confirmationService: ConfirmationService
    ) {}

    openDialog(): void {
        this.displayDialog = true;
    }

    onDialogHide(): void {
        this.confirmationService.confirm({
          message: 'Voulez-vous vraiment quitter sans enregistrer?',
          acceptLabel: 'Oui',
          rejectLabel: 'Non',
          accept: () => {
            this.displayDialog = false;
            this.dialogForm.clearFormAndIndex();
          },
          reject: () => {
            this.displayDialog = true;
          },
        });
    }
}
