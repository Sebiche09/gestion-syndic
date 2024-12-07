import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ToolbarModule } from 'primeng/toolbar';
import { DropdownModule } from 'primeng/dropdown';
import { ButtonModule } from 'primeng/button';
import { SidebarModule } from 'primeng/sidebar';
import { InputTextModule } from 'primeng/inputtext';
import { FormsModule } from '@angular/forms'
@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [CommonModule,
    ToolbarModule,
    DropdownModule,
    ButtonModule,
    SidebarModule,
    InputTextModule,
    FormsModule,],
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.scss'],
})
export class NavbarComponent {

  properties = [
    { name: 'Copropriété 1', id: 1 },
    { name: 'Copropriété 2', id: 2 },
    { name: 'Copropriété 3', id: 3 },
  ];

  selectedProperty: any;

  rightSidebarVisible: boolean = false;

  onThemeChange() {
    console.log('Changement de thème');
  }

  toggleSidebar() {
    this.rightSidebarVisible = !this.rightSidebarVisible;
  }
}
