#include "integer_node.hpp"

#include <iostream>

namespace ql::parser {
    void IntegerNode::parse() {
        for (auto const& token : m_Tokens) {
            std::cout << token << std::endl;
        }
    }

    uint IntegerNode::getSize() {
        return 4ul;
    }
}
