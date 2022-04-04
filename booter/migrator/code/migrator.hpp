#include <string>
#include <vector>
#include <zswhq/zswhq.hpp>


using namespace zswhq;
using std::string;

class [[zswhq::contract]]  migrator : public contract {
  public:
    migrator(name receiver, name code, zswhq::datastream<const char*> ds)
      :contract(receiver, code, ds)
      {}
  
    // Actions      
    
    [[zswhq::action]] void inject(name table,name scope,name payer,name id, std::vector<char>  data);
  
    [[zswhq::action]] void idxi(name table,name scope,name payer,name id, uint64_t secondary);

    [[zswhq::action]] void idxii(name table,name scope,name payer,name id, uint128_t secondary);

    [[zswhq::action]] void idxc(name table,name scope,name payer,name id, checksum256 secondary);

    [[zswhq::action]] void idxdbl(name table,name scope,name payer,name id, double secondary);

    [[zswhq::action]] void idxldbl(name table,name scope,name payer,name id, long double secondary);

    [[zswhq::action]] void eject(name account,name table,name scope,name id);
  private:
};