#pragma once

#include "parser/node/parse_node.hpp"

namespace ql::parser {
    class PackageNode : public ParseNode {
    private:
        std::string m_Name;
    public:
        using ParseNode::ParseNode;

        void parse() override;
    };
}