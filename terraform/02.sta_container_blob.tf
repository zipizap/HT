resource "azurerm_storage_account" "mysta321" {
    name                        = "newsta321${var.project_name}"
    resource_group_name         = azurerm_resource_group.rg.name
    location                    = azurerm_resource_group.rg.location
    account_kind                = "StorageV2"
    account_tier                = "Standard"
    account_replication_type    = "LRS"
    access_tier                 = "Hot"
    allow_blob_public_access    = "true"
}
output "primary_blob_endpoint" {
    value                       = azurerm_storage_account.mysta321.primary_blob_endpoint
}
output "primary_blob_host" {
    value                       = azurerm_storage_account.mysta321.primary_blob_host
}
output "primary_access_key" {
    value                       = azurerm_storage_account.mysta321.primary_access_key
    sensitive = true
}

resource "azurerm_storage_container" "mycontainer" {
  name                  = "mycontainer"
  storage_account_name  = azurerm_storage_account.mysta321.name
  container_access_type = "container"
}


resource "azurerm_storage_blob" "FileInBlob" {
    name                        = "FileInBlob.txt"
    storage_account_name        = azurerm_storage_account.mysta321.name
    storage_container_name      = azurerm_storage_container.mycontainer.name
    type                        = "Block"
    source                      = "mycontainer/FileInBlob.txt"
    content_type                = "text/html; charset=UTF-8"
}

