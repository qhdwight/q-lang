#include <iostream>

#include "impl_func_node.hpp"

namespace ql::parser {
    void ImplementFunctionNode::parse() {
        std::cout << "IMPL: " << m_Body << "#" << std::endl;
    }
}
