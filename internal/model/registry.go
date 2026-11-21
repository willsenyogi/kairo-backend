package model

func AllModels() []interface{} {
    return []interface{}{
        &MasterProvince{},
        &MasterCity{},
        &MasterDistrict{},
        &MasterSubdistrict{},
        &MasterPostalCode{},
        &MasterSubdistrictPostalCode{},

        &User{},
        &Seller{},
        &Category{},
        &Product{},
        &ProductImage{},
        &Address{},
        &OrderHeader{},
        &OrderDetail{},
        &Payment{},
        &AuditLog{},
    }
}