import { Component, ViewChild, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { DialogModule } from 'primeng/dialog';
import { TableModule } from 'primeng/table';
import { ButtonModule } from 'primeng/button';
import { CondominiumComponent } from '../condominium/condominium.component';
import { ConfirmationService } from 'primeng/api';
import { ConfirmDialogModule } from 'primeng/confirmdialog';
import { CondominiumService } from '../../services/condominium.service';
@Component({
    selector: 'displayCondominium',
    templateUrl: 'display-condominium.component.html',
    standalone: true,
    imports: [TableModule, ButtonModule, DialogModule, CondominiumComponent, ConfirmDialogModule],
    providers: [ConfirmationService, CondominiumService]
})
export class DisplayCondominiumComponent implements OnInit {
    displayDialog: boolean = false;
    condominiums: any[] = [];
    @ViewChild(CondominiumComponent) dialogForm!: CondominiumComponent; 

    constructor(
        private router: Router,
        private confirmationService: ConfirmationService,
        private condominiumService: CondominiumService
        
    ) {}

    ngOnInit(): void {
      this.condominiumService.getAllCondominiums().subscribe(
        (data: any) => {
          this.condominiums = data;
          console.log('Données des condominiums:', this.condominiums);
        },
        (error) => {
          console.error('Erreur lors de la récupération des données:', error);
        }
      );
    }
    
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
