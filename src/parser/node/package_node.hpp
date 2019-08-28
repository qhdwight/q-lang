#pragma once

#include "parse_node.hpp"

namespace ql::parser {
    class PackageNode : public ParseNode {
    private:
        std::string m_Name;
    public:
        PackageNode(std::string const& rawText, ParentRef const& parent)
                : ParseNode(rawText, parent) {}

        void parse(std::string const& text) override;
    };
}