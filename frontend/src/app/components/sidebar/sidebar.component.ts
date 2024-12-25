import { Component, ViewChild } from '@angular/core';
import { SidebarModule } from 'primeng/sidebar';
import { ButtonModule } from 'primeng/button';
import { RippleModule } from 'primeng/ripple';
import { AvatarModule } from 'primeng/avatar';
import { StyleClassModule } from 'primeng/styleclass';
import { RouterLink } from '@angular/router';
import { CommonModule } from '@angular/common';
@Component({
  selector: 'app-sidebar',
  standalone: true,
  imports: [SidebarModule, ButtonModule, RippleModule, AvatarModule, StyleClassModule,RouterLink, CommonModule],
  templateUrl: './sidebar.component.html',
  styleUrl: './sidebar.component.scss'
})
export class SidebarComponent {

  sidebarHovered: boolean = false;
  mouseleaveTimeout: any;

  closeSidebar() {
    this.sidebarHovered = false;
  }
  menuItems = [
    { label: 'COPROPRIETE(S)', icon: 'pi pi-building', link: '/displaycondominium', title: 'DASHBOARDS' },
    { label: 'COPROPRIETAIRE(S)', icon: 'pi pi-user' },
    { label: 'OCCUPANT(S)', icon: 'pi pi-user' },
    { label: 'LOT(S)', icon: 'pi pi-home' },
    { label: 'FOURNISSEUR(S)', icon: 'pi pi-truck' },
    { label: 'Reports', icon: 'pi pi-chart-line', title: 'REPORTS' },
    { label: 'Team', icon: 'pi pi-users' },
    { label: 'Messages', icon: 'pi pi-comments', badge: '3' },
    { label: 'Calendar', icon: 'pi pi-calendar' },
    { label: 'Settings', icon: 'pi pi-cog' },
  ];
  onMouseLeave() {
    if (this.mouseleaveTimeout) {
      clearTimeout(this.mouseleaveTimeout);
    }
  
    this.mouseleaveTimeout = setTimeout(() => {
      this.sidebarHovered = false;
    }, 200);
  }

}
