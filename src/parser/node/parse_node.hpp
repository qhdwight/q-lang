#pragma once

#include "abstract_node.hpp"

namespace ql::parser {
    class ParseNode : public AbstractNode {
    protected:
        std::string m_RawText;
        std::string_view m_Body;
        std::vector<std::string> m_Tokens;
    public:
        ParseNode(std::string&& rawText, std::string_view const& body, std::vector<std::string>&& tokens, ParentRef const& parent)
                : AbstractNode(parent), m_Body(body), m_RawText(rawText), m_Tokens(tokens) {
        }

        virtual void parse() {};

        std::string_view getText() const { return m_RawText; }
    };
}
