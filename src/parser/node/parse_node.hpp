#pragma once

#include "abstract_node.hpp"

#include <utility>

namespace ql::parser {
    class ParseNode : public AbstractNode {
    private:
        std::string m_RawText;
        std::vector<std::string> m_Tokens;
    public:
        ParseNode(std::string&& rawText, std::vector<std::string>&& tokens, ParentRef const& parent)
                : AbstractNode(parent), m_RawText(rawText), m_Tokens(tokens) {
            parse(m_RawText, m_Tokens);
        }

        std::string_view getText() const { return m_RawText; }

        virtual void parse(std::string const& text, const std::vector<std::string>& tokens) {};
    };
}
