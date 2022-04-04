#include "migrator.hpp"

#define IMPORT extern "C" __attribute__((zswhq_wasm_import))

extern "C" {
  #include "zswhq/types.h"
}

using namespace zswhq;

IMPORT int32_t db_store_i64(uint64_t scope, uint64_t table, uint64_t payer, uint64_t id,  const void* data, uint32_t len);
IMPORT void db_remove_i64(int32_t iterator);
IMPORT int32_t db_idx64_store(uint64_t scope, capi_name table, capi_name payer, uint64_t id, const uint64_t* secondary);
IMPORT int32_t db_idx128_store(uint64_t scope, capi_name table, capi_name payer, uint64_t id, const uint128_t* secondary);
IMPORT int32_t db_idx256_store(uint64_t scope, capi_name table, capi_name payer, uint64_t id, const uint128_t* data, uint32_t data_len );
IMPORT int32_t db_idx_double_store(uint64_t scope, capi_name table, capi_name payer, uint64_t id, const double* secondary);
IMPORT int32_t db_idx_long_double_store(uint64_t scope, capi_name table, capi_name payer, uint64_t id, const long double* secondary);


void migrator::inject(name table,name scope,name payer,name id, std::vector<char>  data) {            
    zswhq::print("inject ", zswhq::name(table), ":", zswhq::name(scope), " <", zswhq::name(id),":", zswhq::name(payer) ,">\n");
    const auto resp = db_store_i64(
      scope.value,      // The scope where the record will be stored
      table.value,      // The ID/name of the table within the current scope/code context
      payer.value,      // The account that is paying for this storage
      id.value,         // Id of the entry
      (void*)&data[0],  // Record to store
      data.size()       // Size of data
    );
    zswhq::print("inject resp: ", resp , "\n");
};

void migrator::idxi(name table,name scope,name payer,name id, uint64_t secondary) {            
    zswhq::print("idxi ", zswhq::name(table), ":", zswhq::name(scope), " <", zswhq::name(id),":", zswhq::name(payer) ,">\n");
    const auto resp = db_idx64_store(
      scope.value,      // The scope for the secondary index
      table.value,      // The ID/name of the table within the current scope/code context
      payer.value,      // The account that is paying for this storage
      id.value,         // Id of the entry
      &secondary  // Record to store
    );
    zswhq::print("idxi resp: ", resp , "\n");
};

void migrator::idxii(name table,name scope,name payer,name id, uint128_t secondary) {            
    zswhq::print("idxii ", zswhq::name(table), ":", zswhq::name(scope), " <", zswhq::name(id),":", zswhq::name(payer) ,">\n");
    const auto resp = db_idx128_store(
      scope.value,      // The scope for the secondary index
      table.value,      // The ID/name of the table within the current scope/code context
      payer.value,      // The account that is paying for this storage
      id.value,         // Id of the entry
      &secondary      // Record to store
    );  
    zswhq::print("idxii resp: ", resp , "\n");
};

void migrator::idxc(name table,name scope,name payer,name id, checksum256 secondary) {
    zswhq::print("idxc ", zswhq::name(table), ":", zswhq::name(scope), " <", zswhq::name(id),":", zswhq::name(payer) ,">\n");
    const auto ref = secondary.get_array();
    const auto resp = db_idx256_store(
      scope.value,      // The scope for the secondary index
      table.value,      // The ID/name of the table within the current scope/code context
      payer.value,      // The account that is paying for this storage
      id.value,         // Id of the entry
      ref.data(),       // Record to store
      2
    ); 
    zswhq::print("idxc resp: ", resp , "\n"); 
};

void migrator::idxdbl(name table,name scope,name payer,name id, double secondary) {
    zswhq::print("idxdbl ", zswhq::name(table), ":", zswhq::name(scope), " <", zswhq::name(id),":", zswhq::name(payer) ,">\n");
    const auto resp = db_idx_double_store(
      scope.value,      // The scope for the secondary index
      table.value,      // The ID/name of the table within the current scope/code context
      payer.value,      // The account that is paying for this storage
      id.value,         // Id of the entry
      &secondary      // Record to store
    );
    zswhq::print("idxdbl resp: ", resp , "\n"); 
};

void migrator::idxldbl(name table,name scope,name payer,name id, long double secondary) {
    zswhq::print("idxldbl ", zswhq::name(table), ":", zswhq::name(scope), " <", zswhq::name(id),":", zswhq::name(payer) ,">\n");
    const auto resp = db_idx_long_double_store(
      scope.value,      // The scope for the secondary index
      table.value,      // The ID/name of the table within the current scope/code context
      payer.value,      // The account that is paying for this storage
      id.value,         // Id of the entry
      &secondary      // Record to store
    );
    zswhq::print("idxldbl resp: ", resp , "\n");  
};

void migrator::eject(name account,name table,name scope,name id) {
  zswhq::print("delete ", zswhq::name(account), ":", zswhq::name(table), ":", zswhq::name(scope) , " <", zswhq::name(id), ">\n");
	int32_t itr = zswhq::internal_use_do_not_use::db_find_i64(account.value, scope.value, table.value, id.value);
	db_remove_i64(itr);
  zswhq::print("idxldbl resp: itr=", itr , "\n");  
};
