#pragma once

#include <utility>

#include "abstract_node.hpp"

namespace ql::parser {
    class ParseNode : public AbstractNode {
    private:
        std::string m_RawText;
    public:
        ParseNode(std::string&& rawText, ParentRef const& parent)
                : AbstractNode(parent), m_RawText(rawText) {
            parse(m_RawText);
        }

        std::string_view getText() const { return m_RawText; }

        virtual void parse(std::string const& text) {};
    };
}
