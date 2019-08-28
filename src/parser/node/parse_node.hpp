#pragma once

#include "abstract_node.hpp"

namespace ql::parser {
    class ParseNode : public AbstractNode {
    private:
        std::string m_RawText;
    public:
        ParseNode(std::string&& rawText, const ParentRef& parent)
                : AbstractNode(parent), m_RawText(rawText) {}

        std::string_view getText() const { return m_RawText; }
    };
}
