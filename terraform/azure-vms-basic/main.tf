terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~>4.1.0"
    }
  }
}

provider "azurerm" {
  features {}
  resource_provider_registrations = "none"
}

variable "resource_group_name" {
  type = string
}

data "azurerm_resource_group" "example" {
  name = var.resource_group_name
}

variable "vm_count" {
  type    = number
  default = 3
}

resource "azurerm_virtual_network" "example" {
  name                = "vnet"
  resource_group_name = data.azurerm_resource_group.example.name
  location            = data.azurerm_resource_group.example.location

  address_space = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "subnet"
  resource_group_name  = data.azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name

  address_prefixes = ["10.0.0.0/24"]
}

variable "admin_password" {
  type      = string
  sensitive = true
}

resource "azurerm_linux_virtual_machine" "example" {
  count = var.vm_count

  name                = "ubuntu-${count.index}"
  resource_group_name = data.azurerm_resource_group.example.name
  location            = data.azurerm_resource_group.example.location

  network_interface_ids = [azurerm_network_interface.example[count.index].id]
  size                  = "Standard_DS1_v2"

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "ubuntu-24_04-lts"
    sku       = "server"
    version   = "latest"
  }

  admin_username                  = "adminuser"
  admin_password                  = var.admin_password
  disable_password_authentication = false
}

resource "azurerm_network_interface" "example" {
  count = var.vm_count

  name                = "nic-${count.index}"
  resource_group_name = data.azurerm_resource_group.example.name
  location            = data.azurerm_resource_group.example.location

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.example.id
    private_ip_address_allocation = "Dynamic"
  }
}