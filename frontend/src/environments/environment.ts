export const environment = {
    production: false,
    apiUrls: {
      uploadApi: 'http://localhost:8080/upload',
      condominiumApi: 'http://localhost:8080/condominium',
      getAllCondominiums: 'http://localhost:8080/all-condominiums',
      getCivilities: 'http://localhost:8080/getcivilities',
      getReminderRemindingMethod : 'http://localhost:8080/getreminderreceivingmethods',
      getDocumentRemindingMethod : 'http://localhost:8080/getdocumentreceivingmethods',
      getCountries: 'https://restcountries.com/v3.1/independent?status=true&fields=name,cca2',
      getCitiesBase: 'http://api.geonames.org/searchJSON',
      username: 'sebastien.vlmnckx',
      checkUniqueness: 'http://localhost:8080/check-uniqueness',
      unit : 'http://localhost:8080/unit'
    }
  };  